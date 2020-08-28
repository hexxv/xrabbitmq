// Copyright 2020/8 @Author:hex
//
// Routing模式(路由模式，一个消息被多个消费者获取，并且消息的目标队列可被生产者指定)
// 1. 消息生产者将消息发送给交换机按照路由判断,路由是字符串(info) 当前产生的消息携带
//    路由字符(对象的方法),交换机根据路由的key,只能匹配上路由key对应的消息队列,对应
//    的消费者才能消费消息;
// 2. 根据业务功能定义路由字符串
// 3. 从系统的代码逻辑中获取对应的功能字符串,将消息任务扔到对应的队列中业务场景:error通知
//    EXCEPTION;错误通知的功能;传统意义的错误通知;客户通知;利用key路由,可以将程序中
//    的错误封装成消息传入到消息队列中,开发者可以自定义消费者,实时接收错误;

package routing

import (
	"xrabbitmq/pkg/consumer"
	"xrabbitmq/pkg/external"
	"xrabbitmq/pkg/log"
	"xrabbitmq/pkg/session"
)

type routing struct {
	*consumer.Consumer
}

func New(sess *session.Session) *routing {
	return &routing{consumer.NewConsumer(sess, consumer.ModelRouting)}
}

func (c *routing) Consume(handler func(delivery external.XDelivery)) (err error) {
	defer c.Done(err)

	queueOptions := c.Sess().Queue()
	consumerOptions := c.Sess().OptionsConsumer()
	exchangeOptions := c.Sess().Exchange()
	BindingOptions := c.Sess().Binding()

	err = c.Sess().Channel().ExchangeDeclare(
		exchangeOptions.Name,
		exchangeOptions.Typ,
		exchangeOptions.Durable,
		exchangeOptions.AutoDelete,
		exchangeOptions.Internal,
		exchangeOptions.NoWait,
		exchangeOptions.Args,
	)
	if err != nil {
		log.Logger.Errorf("%s.ExchangeDeclare error: %s", c.Model(), err)
		return err
	}

	q, err := c.Sess().Channel().QueueDeclare(
		queueOptions.Name, // 随机生产队列名称,这里注意队列名称不要写
		queueOptions.Durable,
		queueOptions.AutoDelete,
		queueOptions.Exclusive,
		queueOptions.NoWait,
		queueOptions.Args,
	)
	if err != nil {
		log.Logger.Errorf("%s.QueueDeclare error: %s", c.Model(), err)
		return err
	}

	err = c.Sess().Channel().QueueBind(
		q.Name,
		BindingOptions.RoutingKey,
		exchangeOptions.Name,
		BindingOptions.NoWait,
		BindingOptions.Args,
	)
	if err != nil {
		log.Logger.Errorf("%s.QueueBind error: %s", c.Model(), err)
		return err
	}

	deliveries, err := c.Sess().Channel().Consume(
		q.Name,
		consumerOptions.Tag,
		consumerOptions.AutoAck,
		consumerOptions.NoLocal,
		consumerOptions.Exclusive,
		consumerOptions.NoWait,
		consumerOptions.Args,
	)
	if err != nil {
		log.Logger.Errorf("%s.Consume error: %s.", c.Model(), err)
		return err
	}

	c.Consumer.Consume(deliveries, handler)
	return nil
}
