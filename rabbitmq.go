package xrabbitmq

import (
	"fmt"
	"sync"
	"xrabbitmq/build"
	"xrabbitmq/pkg/external"
	"xrabbitmq/pkg/log"
	"xrabbitmq/pkg/session"
)

// RabbitMQ客户端
type RabbitMQ struct {
	// rabbitMQ 客户端连接
	conn *external.XConnection

	// rabbitMq 通信管道
	channel *external.XChannel

	// 建立连接需要的配置项
	ConfOptions

	// 启动
	startupOnce sync.Once

	// 关闭
	shutdownOnce sync.Once

	// 关闭通知
	closing chan struct{}

	wg sync.WaitGroup
}

// New 根据配置项初始化并返回一个RabbitMQ客户端实例
func New(opts ...ConfOption) *RabbitMQ {
	return &RabbitMQ{
		closing:     make(chan struct{}),
		ConfOptions: defaultConfOptions(opts...),
	}
}

// Conn 得到RabbitMQ客户端与目标RabbitMQ服务端之间所建立的连接
func (rmq *RabbitMQ) Conn() *external.XConnection {
	return rmq.conn
}

// Startup 启动RabbitMQ客户端
// 也就是RabbitMQ客户端与目标RabbitMQ服务端在此时会建立连接(仅一次)
func (rmq *RabbitMQ) Startup() error {
	var err error
	rmq.startupOnce.Do(func() {
		if rmq.conn != nil {
			return
		}
		err = rmq.dial()
	})
	return err
}

// Shutdown 关闭RabbitMQ客户端连接
// 一般用于程序退出前释放当前客户端与服务端之间建立的连接
func (rmq *RabbitMQ) Shutdown() error {
	log.Logger.Warning("RabbitMQ will shutdown...")
	var err error
	rmq.shutdownOnce.Do(func() {
		close(rmq.closing)
		rmq.wg.Wait()
		if err = rmq.conn.Close(); err != nil {
			if xerr, ok := err.(*external.XError); ok {
				if xerr.Code != external.ChannelError {
					err = fmt.Errorf("AMQP connection close error: %w", err)
					log.Logger.Error("RabbitMQ shutdown error: %s", err)
					return
				}
			}
		}
		log.Logger.Info("RabbitMQ shutdown OK")
	})
	return err
}

// BuildConsumer 得到消费者构建工具
func (rmq *RabbitMQ) BuildConsumer(opts ...session.Option) external.ConsumerBuilder {
	return build.NewConsumerBuild(
		build.DependConn(rmq.conn),
		opts...,
	)
}

// BuildProducer 得到生产者构建工具
func (rmq *RabbitMQ) BuildProducer(opts ...session.Option) external.ProducerBuilder {
	return build.NewProducerBuild(
		build.DependConn(rmq.conn),
		opts...,
	)
}

// dial 顾名思义
func (rmq *RabbitMQ) dial() error {
	conf := external.XURI{
		Scheme:   "amqp",
		Host:     rmq.Host,
		Port:     rmq.Port,
		Username: rmq.User,
		Password: rmq.Pwd,
		Vhost:    rmq.VHost,
	}.String()

	var err error
	rmq.conn, err = external.XDial(conf)
	if err != nil {
		return err
	}
	rmq.channel, err = rmq.conn.Channel()
	if err != nil {
		return err
	}

	rmq.wg.Add(1)
	go rmq.handleErrors(rmq.conn)

	return nil
}

// handleErrors 开启一个goroutine来处理/监听连接过程中所发生的错误
func (rmq *RabbitMQ) handleErrors(conn *external.XConnection) {
	defer func() {
		if x := recover(); x != nil {
			log.Logger.Errorf("Panic:%+v", x)
		}
		rmq.wg.Done()
	}()

	closeChan := conn.NotifyClose(make(chan *external.XError))
	blockChan := conn.NotifyBlocked(make(chan external.XBlocking))

	for {
		select {
		case xerr := <-closeChan:
			// CRITICAL Exception (320) Reason: "CONNECTION_FORCED - broker forced connection closure with reason 'shutdown'"
			// CRITICAL Exception (501) Reason: "read tcp 127.0.0.1:5672: i/o timeout"
			// CRITICAL Exception (503) Reason: "COMMAND_INVALID - unimplemented method"
			switch xerr.Code {
			case external.NotFound: // 404
				log.Logger.Error("amqp.NotFound")
			case external.FrameError: // 501
				log.Logger.Error("amqp.FrameError")
			case external.ConnectionForced: // 320
				log.Logger.Error("amqp.ConnectionForced")
			}
		case b := <-blockChan:
			if b.Active {
				log.Logger.Error("TCP blocked: %q", b.Reason)
			} else {
				log.Logger.Error("TCP unblocked")
			}
		case _, ok := <-rmq.closing:
			if !ok {
				log.Logger.Warning("handleErrors recv rmq shutdown, return now.")
				return
			}
		}
	}
}
