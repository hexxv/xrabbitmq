// Copyright 2020/8 @Author:hex
//
// work 工作模式
// 1. 消息产生者将消息放入队列消费者可以有多个,消费者1,消费者2,同时监听同一个队列,消息被消费?
//    C1 C2共同争抢当前的消息队列内容,谁先拿到谁负责消费消息(隐患,高并发情况下,默认会产生某一个
//    消息被多个消费者共同使用,可以设置一个开关(syncronize,与同步锁的性能不一样) 保证一条消息
//    只能被一个消费者使用)
// 2. 应用场景:红包;大项目中的资源调度(任务分配系统不需知道哪一个任务执行系统在空闲,直接将任务扔
//    到消息队列中,空闲的系统自动争抢)

package work

import (
	"xrabbitmq/pkg/consumer"
	"xrabbitmq/pkg/external"
	"xrabbitmq/pkg/log"
	"xrabbitmq/pkg/session"
)

type work struct {
	*consumer.Consumer
}

func New(sess *session.Session) *work {
	return &work{consumer.NewConsumer(sess, consumer.ModelWork)}
}

// 开始消费
func (c *work) Consume(handler func(delivery external.XDelivery)) (err error) {
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
