// Copyright 2020/8 @Author:hex
//
// Publish模式（订阅模式，消息被路由投递给多个队列，一个消息被多个消费者获取）
// 1. X代表交换机rabbitMQ内部组件,erlang 消息产生者是代码完成,代码的执行效率不高, 消息产生者
//    将消息放入交换机,交换机发布订阅把消息发送到所有消息队列中,对应消息队列的消费者拿到消息进行
//    消费
// 2. 应用场景:邮件群发,群聊天,广播(广告)

package publish

import (
	"xrabbitmq/pkg/consumer"
	"xrabbitmq/pkg/external"
	"xrabbitmq/pkg/log"
	"xrabbitmq/pkg/session"
)

type publish struct {
	*consumer.Consumer
}

func New(sess *session.Session) *publish {
	return &publish{consumer.NewConsumer(sess, consumer.ModelPublish)}
}

func (c *publish) Consume(handler func(delivery external.XDelivery)) (err error) {
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
