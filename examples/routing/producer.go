package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"xrabbitmq"
	"xrabbitmq/pkg/external"
	"xrabbitmq/pkg/session"
	"xrabbitmq/pkg/session/broker"
	"xrabbitmq/pkg/session/broker/binding"
	"xrabbitmq/pkg/session/broker/exchange"
	"syscall"
	"time"
)

var pid int

func init() {
	flag.IntVar(&pid, "p", 0, "生产者id")
}

func main() {

	flag.Parse()

	if pid <= 0 {
		fmt.Println("use -p to give the producer an ID")
		return
	}

	// rabbitMQ实例
	rabbitMQ := xrabbitmq.New(
		xrabbitmq.WithHost("localhost"),
		xrabbitmq.WithPort(5672),
		xrabbitmq.WithUser("hexx"),
		xrabbitmq.WithPwd("hexx"),
		xrabbitmq.WithVHost("/"),
	)

	// 启动rabbitMQ建立连接
	err := rabbitMQ.Startup()
	if err != nil {
		fmt.Println("mq startup error:", err)
		return
	}
	// 程序退出时，断掉连接，释放资源
	defer rabbitMQ.Shutdown()

	// 得到一个Simple模式的生产者
	producer, err := rabbitMQ.BuildProducer(
		// 设置broker中exchange的name
		session.WithBrokerOptions(
			broker.WithExchange(exchange.SetName("hexT_routing_testHelloWordE?change1")),
			broker.WithBinding(binding.SetRoutingKey(fmt.Sprintf("hex_routingKey_%d", pid))),
		),
	).Routing(false)

	if err != nil {
		fmt.Println("get mq producer error:", err)
		return
	}
	// 使用完后，关掉该生产者所占用的通信管道(通信管带来自于rabbitMQ的连接,至于为什么？去网上学习吧)
	defer producer.Cancel()

	// 开启一个协程，监听未成功路由的消息
	// producer.NotifyReturn(func(xReturn external.XReturn) {
	// 	fmt.Println("NotifyReturn:", xReturn)
	// })
	messages := make(chan *external.XPublishMsg, 5000)
	done := make(chan error)
	go func() {
		err = producer.Publish(messages)
		if err != nil {
			done <- err
			return
		}
	}()

	fmt.Println("producer running, please press `Ctrl+C` to stop.")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)

	ticker := time.NewTicker(time.Second * 5)
	defer func() {
		ticker.Stop()
		close(messages)
		fmt.Println("Stopped.")
	}()

	var i int
	for {
		// 监听Ctrl+C及生产者生产消息时发生的错误
		select {
		case <-ticker.C:
			i++
			msg := fmt.Sprintf("No.%d producer say hello my name is hex%d", pid, i)
			fmt.Println(msg)
			messages <- &external.XPublishMsg{Body: []byte(msg)}
		case sig := <-signals:
			fmt.Println("handle signal: ", sig)
			return
		case err := <-done:
			fmt.Println("producer error: ", err)
			return
		}
	}

}
