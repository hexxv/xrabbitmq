// Copyright 2020/8 @Author:hex
//
// Simple模式
// 生产环境中不推荐使用Simple模式, 执行过程及不推荐的原因如下:
// 1.消息产生者将消息放入队列
// 2.消息的消费者(consumer) 监听(while) 消息队列,如果队列中有消息,就消费掉。
// 3.消息被拿走后,自动从队列中删除(隐患 消息可能没有被消费者正确处理,已经从队列中消失了,造成消息的丢失)

package simple

import (
	"xrabbitmq/pkg/consumer"
	"xrabbitmq/pkg/external"
	"xrabbitmq/pkg/log"
	"xrabbitmq/pkg/session"
)

type simple struct {
	*consumer.Consumer
}

func New(sess *session.Session) *simple {
	return &simple{consumer.NewConsumer(sess, consumer.ModelSimple)}
}

// 开始消费
func (c *simple) Consume(handler func(delivery external.XDelivery)) (err error) {
	defer c.Done(err)

	queueOptions := c.Sess().Queue()
	consumerOptions := c.Sess().OptionsConsumer()
	q, err := c.Sess().Channel().QueueDeclare(
		queueOptions.Name,
		queueOptions.Durable,
		queueOptions.AutoDelete,
		queueOptions.Exclusive,
		queueOptions.NoWait,
		queueOptions.Args,
	)
	if err != nil {
		log.Logger.Errorf("%s.QueueDeclare error: %s.", c.Model(), err)
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
