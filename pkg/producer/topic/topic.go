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
	"xrabbitmq/pkg/external"
	"xrabbitmq/pkg/log"
	"xrabbitmq/pkg/producer"
	"xrabbitmq/pkg/session"
)

type topic struct {
	*producer.Producer
}

func New(sess *session.Session) *topic {
	return &topic{producer.NewProducer(sess, producer.ModelTopic)}
}

func NewDynamic(sess *session.Session) *topic {
	return &topic{producer.NewProducer(sess, producer.ModelTopicDynamic)}
}

func (p *topic) Publish(messages <-chan *external.XPublishMsg) (err error) {
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
