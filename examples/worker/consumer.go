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
	"syscall"
)

var cid int

func init() {
	flag.IntVar(&cid, "c", 0, "消费者id")
}

func main() {
	flag.Parse()

	if cid <= 0 {
		fmt.Println("use -p to give the consumer an ID")
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

	// 得到一个Simple模式的消费者
	consumer, err := rabbitMQ.BuildConsumer(
		// 设置broker中queue的name 要与生产者对应，不然不在一个队列
		session.WithBrokerOptions(
			broker.WithQueue(queue.SetName("hexTestRabbitMQ1_HelloWord!")),
		),
	).Work()
	if err != nil {
		fmt.Println("get mq producer error:", err)
		return
	}
	// 使用完后，关掉该消费者所占用的通信管道
	defer consumer.Cancel()

	err = consumer.Qos(1)
	if err != nil {
		fmt.Println("consumer.Qos error: ", err)
		return
	}

	fmt.Println("try to consumer...")
	done := make(chan error)

	go func() {
		err = consumer.Consume(ConsumerHandler)
		if err != nil {
			done <- err
			return
		}
	}()

	fmt.Println("consumer running, please press `Ctrl+C` to stop.")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)

	select {
	case sig := <-signals:
		fmt.Println("handle signal: ", sig)
	case err := <-done:
		fmt.Println("consumer error: ", err)
	}

	fmt.Println("Stopped.")
}

func ConsumerHandler(delivery external.XDelivery) {
	fmt.Println(fmt.Sprintf("No.%d consumer recv:%s \n", cid, string(delivery.Body)))
	delivery.Ack(false)
}
