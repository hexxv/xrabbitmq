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
	"xrabbitmq/pkg/session/broker/queue"
	"xrabbitmq/pkg/session/produceropts"
	"sync"
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
		// 设置broker中queue的name
		session.WithBrokerOptions(
			broker.WithQueue(queue.SetName("hexTestRabbitMQ1_HelloWord!")),
		),
		// 监听消息未被路由到(可以简单理解为未送达)的消息的回调
		session.WithPublishingOptions(
			produceropts.SetMandatory(true),
		),
	).Work()

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
	messages := make(chan *external.XPublishMsg, 1000000)
	done := make(chan error)
	go func() {
		err = producer.Publish(messages)
		if err != nil {
			done <- err
			return
		}
	}()

	wg := sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var j int
			ticker := time.NewTicker(time.Millisecond * 200)
			for {
				<-ticker.C
				j++
				if j >= 100 {
					return
				}
				msg := fmt.Sprintf("No.%d(%d) producer say hello my name is hex%d", pid, i, j)
				fmt.Println(msg)
				messages <- &external.XPublishMsg{
					Body: []byte(msg),
				}
			}
		}(i)
	}

	wg.Wait()

	fmt.Println("producer running, please press `Ctrl+C` to stop.")
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)

	defer func() {
		close(messages)
		fmt.Println("Stopped.")
	}()

	// 监听Ctrl+C及生产者生产消息时发生的错误
	select {
	case sig := <-signals:
		fmt.Println("handle signal: ", sig)
		return
	case err := <-done:
		fmt.Println("producer error: ", err)
		return
	}

}
