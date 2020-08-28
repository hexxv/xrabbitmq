package main

import (
	"fmt"
	"xrabbitmq"
	"xrabbitmq/pkg/external"
)

func main() {
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

	channel, err := rabbitMQ.Conn().Channel()
	if err != nil {
		fmt.Println("mq channel error:", err)
		return
	}

	exchange1 := "hex.test.1q2ne.e1"
	exchange2 := "hex.test.1q2ne.e2"

	err = channel.ExchangeDeclare(
		exchange1,
		external.XExchangeTopic,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println("mq ExchangeDeclare1 error:", err)
		return
	}

	err = channel.ExchangeDeclare(
		exchange2,
		external.XExchangeTopic,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println("mq ExchangeDeclare2 error:", err)
		return
	}

	q, err := channel.QueueDeclare(
		"hex.test.1q2ne.q",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println("mq QueueDeclare error:", err)
		return
	}

	q, err = channel.QueueDeclare(
		"hex.test.1q2ne.q",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println("mq QueueDeclare2 error:", err)
		return
	}

	err = channel.QueueBind(q.Name, "", exchange1, false, nil)
	if err != nil {
		fmt.Println("mq QueueBind1 error:", err)
		return
	}

	err = channel.QueueBind(q.Name, "", exchange2, false, nil)

	if err != nil {
		fmt.Println("mq QueueBind2 error:", err)
		return
	}

	recv, err := channel.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		fmt.Println("mq Consume error:", err)
		return
	}
	fmt.Println("consume...")
	for msg := range recv {
		msg.Ack(false)
		data := string(msg.Body)
		fmt.Println("recv:", data)
		if data == "ce1" {
			fmt.Println("QueueUnbind1")
			err = channel.QueueUnbind(q.Name, "", exchange1, nil)
			if err != nil {
				fmt.Println("mq QueueUnbind1 error:", err)
				return
			}
		}
		if data == "be1" {
			fmt.Println("QueueBind11")
			err = channel.QueueBind(q.Name, "", exchange1, false, nil)

			if err != nil {
				fmt.Println("mq QueueBind11 error:", err)
				return
			}
		}

	}
}
