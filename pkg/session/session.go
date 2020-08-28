// 生产者/消费者与RabbitMQ之间抽象的一个会话层
package session

import (
	"fmt"
	"xrabbitmq/pkg/external"
	"xrabbitmq/pkg/session/broker"
	"xrabbitmq/pkg/session/broker/binding"
	"xrabbitmq/pkg/session/broker/exchange"
	"xrabbitmq/pkg/session/broker/queue"
	"xrabbitmq/pkg/session/consumeropts"
	"xrabbitmq/pkg/session/produceropts"
)

type Option func(session *Session)

// Session 是一个较为抽象的概念，它保存了当前生产者/消费者所占用的通信管道。同时,
// 包含了当前会话所携带的诸多属性：包括交换机、队列、binding、生产者/消费者，所以
// 可以抽象的认定为是当前生产者/消费者与RabbitMQ之间的一个会话，虽然RabbitMQ的
// 架构模型中并没有这层概念
type Session struct {
	// channel 生产者/消费者基于上层与RabbitMQ建立的连接(amqp.Connection)在此连接上开辟出来的一条通信管道
	// 对于操作系统而言，建立连接是很消耗资源的。相比而言，基于一条连接开辟多条通信管道是更加高效、轻量的方式
	channel *external.XChannel

	// Broker 是RabbitMQ架构模型中的中介模块
	// 它包含了交换机、队列、binding。 更多关于broker的知识可以上网查询
	broker broker.Broker

	// ConsumerOptions 消费者的配置项
	consumerOptions consumeropts.Options

	// PublishingOptions 生产者的配置项
	producerOptions produceropts.Options
}

// Establish 建立通信管道
func (s *Session) Establish(conn *external.XConnection) error {
	channel, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("session establish error: get channel by connection error: %s", err)
	}
	s.channel = channel
	return nil
}

// HasEstablished 是否建立了通信管道
func (s *Session) HasEstablished() bool {
	return s.channel != nil
}

// NewSession 传入各种配置项，得到一个未建立通信管道的会话实例
// 在使用中会根据不用模型下的生产者/消费者构建需求，提前对这些配置项做校验
// 在校验通过后，调用 session.Establish(conn rmq连接)来建立通信管道
// 最终得到一个完整的session
func NewSession(opts ...Option) *Session {
	s := Session{
		channel:         nil,
		broker:          broker.Broker{},
		consumerOptions: consumeropts.Options{},
		producerOptions: produceropts.Options{},
	}
	for _, o := range opts {
		o(&s)
	}
	return &s
}

// Channel get current session‘s channel
func (s *Session) Channel() *external.XChannel {
	return s.channel
}

// Exchange get exchange setting
func (s *Session) Exchange() exchange.Exchange {
	return s.broker.Exchange
}

// Queue get queue setting
func (s *Session) Queue() queue.Queue {
	return s.broker.Queue
}

// Binding get binding setting
func (s *Session) Binding() binding.Binding {
	return s.broker.Binding
}

// OptionsConsumer get consumerOptions
func (s *Session) OptionsConsumer() consumeropts.Options {
	return s.consumerOptions
}

// OptionsProducer get producerOptions
func (s *Session) OptionsProducer() produceropts.Options {
	return s.producerOptions
}

func WithBrokerOptions(bos ...broker.Option) Option {
	// TODO 这里也许可以设置一些默认值？
	return func(session *Session) {
		// bo := broker.Broker{
		// 	Exchange: exchange.Exchange{
		// 		Name:       "",
		// 		Typ:        "",
		// 		Durable:    false,
		// 		AutoDelete: false,
		// 		Internal:   false,
		// 		NoWait:     false,
		// 		Args:       nil,
		// 	},
		// 	Queue: queue.Queue{
		// 		Name:       "",
		// 		Durable:    false,
		// 		AutoDelete: false,
		// 		Exclusive:  false,
		// 		NoWait:     false,
		// 		Args:       nil,
		// 	},
		// 	Binding: binding.Binding{
		// 		RoutingKey: "",
		// 		NoWait:     false,
		// 		Args:       nil,
		// 	},
		// }

		for _, o := range bos {
			o(&session.broker)
		}
		// session.broker = bo
	}
}

func WithConsumerOptions(cos ...consumeropts.Option) Option {
	return func(session *Session) {
		// co := consumeropts.Options{
		// 	Tag:       "",
		// 	AutoAck:   false,
		// 	Exclusive: false,
		// 	NoLocal:   false,
		// 	NoWait:    false,
		// 	Args:      nil,
		// }
		for _, o := range cos {
			o(&session.consumerOptions)
		}
		// session.consumerOptions = co
	}
}

func WithPublishingOptions(pos ...produceropts.Option) Option {
	return func(session *Session) {
		// po := produceropts.Options{
		// 	RoutingKey: "",
		// 	Tag:        "",
		// 	Mandatory:  false,
		// }
		for _, o := range pos {
			o(&session.producerOptions)
		}
		// session.producerOptions = po
	}
}
