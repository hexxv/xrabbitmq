// Copyright 2020/8 @Author:hex
//
// Topic模式（话题模式，一个消息被多个消费者获取）
// 消息的目标queue可用BindingKey以通配符，（#：一个或多个词，*：一个词）的方式指定
// 1. 星号井号代表通配符
// 2. 星号代表多个单词,井号代表一个单词
// 3. 路由功能添加模糊匹配
// 4. 消息产生者产生消息,把消息交给交换机
// 5. 交换机根据key的规则模糊匹配到对应的队列,由队列的监听消费者接收消息消费

package topic

import (
	"xrabbitmq/pkg/consumer"
	"xrabbitmq/pkg/external"
	"xrabbitmq/pkg/log"
	"xrabbitmq/pkg/session"
)

type topic struct {
	*consumer.Consumer
}

func New(sess *session.Session) *topic {
	return &topic{consumer.NewConsumer(sess, consumer.ModelTopic)}
}

func (c *topic) Consume(handler func(delivery external.XDelivery)) (err error) {
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
