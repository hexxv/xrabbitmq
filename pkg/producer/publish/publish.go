// Copyright 2020/8 @Author:hex
//
// Publish模式（订阅模式，消息被路由投递给多个队列，一个消息被多个消费者获取）
// 1. X代表交换机rabbitMQ内部组件,erlang 消息产生者是代码完成,代码的执行效率不高, 消息产生者
//    将消息放入交换机,交换机发布订阅把消息发送到所有消息队列中,对应消息队列的消费者拿到消息进行
//    消费
// 2. 应用场景:邮件群发,群聊天,广播(广告)

package publish

import (
	"xrabbitmq/pkg/external"
	"xrabbitmq/pkg/log"
	"xrabbitmq/pkg/producer"
	"xrabbitmq/pkg/session"
)

type publish struct {
	*producer.Producer
}

func New(sess *session.Session) *publish {
	return &publish{producer.NewProducer(sess, producer.ModelPublish)}
}

func (p *publish) Publish(messages <-chan *external.XPublishMsg) (err error) {
	defer p.Done(err)
	// queueOptions := p.Sess().Queue()
	exchangeOptions := p.Sess().Exchange()

	err = p.Sess().Channel().ExchangeDeclare(
		exchangeOptions.Name,
		exchangeOptions.Typ,
		exchangeOptions.Durable,
		exchangeOptions.AutoDelete,
		exchangeOptions.Internal,
		exchangeOptions.NoWait,
		exchangeOptions.Args,
	)

	if err != nil {
		log.Logger.Errorf("%s.ExchangeDeclare error: %s", p.Model(), err)
		return err
	}

	p.Producer.Publish(messages)

	return nil
}
