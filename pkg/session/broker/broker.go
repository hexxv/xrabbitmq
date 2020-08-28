package broker

import (
	"xrabbitmq/pkg/session/broker/binding"
	"xrabbitmq/pkg/session/broker/exchange"
	"xrabbitmq/pkg/session/broker/queue"
)

type Option func(broker *Broker)

// RabbitMQ Broker[中介]
type Broker struct {
	// 交换机
	exchange.Exchange
	// 队列
	queue.Queue
	// 绑定
	binding.Binding
}

func WithExchange(eos ...exchange.Option) Option {
	return func(broker *Broker) {
		// eo := exchange.Exchange{
		// 	Name:       "",
		// 	Typ:        "",
		// 	Durable:    false,
		// 	AutoDelete: false,
		// 	Internal:   false,
		// 	NoWait:     false,
		// 	Args:       nil,
		// }
		for _, o := range eos {
			o(&broker.Exchange)
		}
		// broker.Exchange = eo
	}
}

func WithQueue(qos ...queue.Option) Option {
	return func(broker *Broker) {
		// qo := queue.Queue{
		// 	Name:       "",
		// 	Durable:    false,
		// 	AutoDelete: false,
		// 	Exclusive:  false,
		// 	NoWait:     false,
		// 	Args:       nil,
		// }
		for _, o := range qos {
			o(&broker.Queue)
		}
		// broker.Queue = qo
	}
}

func WithBinding(bos ...binding.Option) Option {
	return func(broker *Broker) {
		// bo := binding.Binding{
		// 	RoutingKey: "",
		// 	NoWait:     false,
		// 	Args:       nil,
		// }
		for _, o := range bos {
			o(&broker.Binding)
		}
		// broker.Binding = bo
	}
}
