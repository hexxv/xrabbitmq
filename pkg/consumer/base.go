package consumer

import (
	"xrabbitmq/pkg/external"
	"xrabbitmq/pkg/internal/utils"
	"xrabbitmq/pkg/log"
	"xrabbitmq/pkg/session"
	"sync"
)

type Consumer struct {
	// session 消费者与交换机/队列/binding的声明保障
	// 消费者会占用一个RabbitMQ的通信管道，不要忘记使用 Consumer.Cancel() 来释放它
	session *session.Session

	// cancelOnce 关闭通信管道幂等
	cancelOnce sync.Once

	// deliveries all deliveries from server will send to this channel
	deliveries <-chan external.XDelivery

	// done: a notifiyng channel for publishings
	// will be used for sync. between close channel and consume handler
	done chan error

	// consumer model simple/work/publish/routing/topic
	model Model
}

func NewConsumer(sess *session.Session, mod Model) *Consumer {
	return &Consumer{
		session: sess,
		done:    make(chan error),
		model:   mod,
	}
}

func (c *Consumer) Sess() *session.Session {
	return c.session
}

func (c *Consumer) Done(err error) {
	c.done <- err
}

func (c *Consumer) Model() Model {
	return c.model
}

func (c *Consumer) Consume(d <-chan external.XDelivery, handler func(delivery external.XDelivery)) {
	c.deliveries = d

	log.Logger.Info("consumer.Consume.handler: deliveries channel starting...")

	// handle all consumer errors, if required re-connect
	// there are problems with reconnection logic for now
	for delivery := range c.deliveries {
		handler(delivery)
	}

	log.Logger.Info("consumer.Consume.handler: deliveries channel closed...")
}

func (c *Consumer) Cancel() error {
	return c.releaseChannel()
}

// Consumer 在断掉连接后 处理之前收到的未处理完的消息
func (c *Consumer) releaseChannel() error {
	var err error
	c.cancelOnce.Do(func() {
		err = utils.CancelChannel(c.session.Channel(), c.session.OptionsConsumer().Tag)
	})
	if err != nil {
		return err
	}
	if c.deliveries == nil {
		close(c.done)
	}
	return <-c.done
}

// Qos controls how many messages the server will try to keep on the network for
// consumers before receiving delivery acks.  The intent of Qos is to make sure
// the network buffers stay full between the server and client.
func (c *Consumer) Qos(messageCount int) error {
	// prefetchCount：消费者未确认消息的个数。
	// prefetchSize ：消费者未确认消息的大小。
	// global ：是否全局生效，true表示是。全局生效指的是针对当前connect里的所有channel都生效。
	return c.session.Channel().Qos(messageCount, 0, false)
}
