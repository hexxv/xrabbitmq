package build

import (
	"fmt"
	"xrabbitmq/pkg/consumer/publish"
	"xrabbitmq/pkg/consumer/routing"
	"xrabbitmq/pkg/consumer/simple"
	"xrabbitmq/pkg/consumer/topic"
	"xrabbitmq/pkg/consumer/work"
	"xrabbitmq/pkg/external"
	"xrabbitmq/pkg/session"
	"xrabbitmq/pkg/session/broker"
	"xrabbitmq/pkg/session/broker/binding"
	"xrabbitmq/pkg/session/broker/exchange"
	"xrabbitmq/pkg/session/broker/queue"
)

type consumerBuild struct {
	sessionOptions []session.Option
	buildRequired  Required
}

func NewConsumerBuild(br Required, sos ...session.Option) *consumerBuild {
	return &consumerBuild{
		sessionOptions: sos,
		buildRequired:  br,
	}
}

func (cb *consumerBuild) sess() *session.Session {
	return session.NewSession(cb.sessionOptions...)
}

func (cb *consumerBuild) rabbitMqConn() *external.XConnection {
	builder := depend{}
	cb.buildRequired(&builder)
	return builder.conn
}

func (cb *consumerBuild) Simple() (external.Consumer, error) {
	cb.sessionOptions = append(cb.sessionOptions, session.WithBrokerOptions(
		broker.WithBinding(
			binding.SetRoutingKey(""),
		),
		broker.WithExchange(
			exchange.SetName(""),
			exchange.SetType(external.XExchangeDirect),
		),
	))
	sess := cb.sess()
	// 如果未指定队列名字，这里不允许由RabbitMQ生成
	if sess.Queue().Name == "" {
		return nil, fmt.Errorf("consumerBuild Simple error: the \"queue's Name\" must be specified")
	}
	err := sess.Establish(cb.rabbitMqConn())
	if err != nil {
		return nil, fmt.Errorf("consumerBuild Simple error: %w", err)
	}
	return simple.New(sess), nil
}

func (cb *consumerBuild) Work() (external.Consumer, error) {
	cb.sessionOptions = append(cb.sessionOptions, session.WithBrokerOptions(
		broker.WithBinding(
			binding.SetRoutingKey(""),
		),
		broker.WithExchange(
			exchange.SetName(""),
			exchange.SetType(external.XExchangeDirect),
		),
	))
	sess := cb.sess()
	// 如果未指定队列名字，这里不允许由RabbitMQ生成
	if sess.Queue().Name == "" {
		return nil, fmt.Errorf("consumerBuild Work error: the \"queue's Name\" must be specified")
	}
	err := sess.Establish(cb.rabbitMqConn())
	if err != nil {
		return nil, fmt.Errorf("consumerBuild Work error: %w", err)
	}
	return work.New(sess), nil
}

func (cb *consumerBuild) Publish() (external.Consumer, error) {
	cb.sessionOptions = append(cb.sessionOptions, session.WithBrokerOptions(
		broker.WithQueue(queue.SetName(""),
			queue.SetExclusive(true),
		),
		broker.WithExchange(
			exchange.SetDurable(true),
			exchange.SetType(external.XExchangeFanout),
		),
		broker.WithBinding(
			binding.SetRoutingKey(""),
		),
	))
	sess := cb.sess()

	if sess.Exchange().Name == "" {
		return nil, fmt.Errorf("consumerBuild Publish error: the \"exchange's Name\" must be specified")
	}

	err := sess.Establish(cb.rabbitMqConn())
	if err != nil {
		return nil, fmt.Errorf("consumerBuild Publish error: %w", err)
	}
	return publish.New(sess), nil
}

func (cb *consumerBuild) Routing() (external.Consumer, error) {
	cb.sessionOptions = append(cb.sessionOptions, session.WithBrokerOptions(
		broker.WithQueue(queue.SetName(""),
			queue.SetExclusive(true),
		),
		broker.WithExchange(
			exchange.SetDurable(true),
			exchange.SetType(external.XExchangeDirect),
		),
	))
	sess := cb.sess()
	if sess.Exchange().Name == "" {
		return nil, fmt.Errorf("consumerBuild Routing error: the \"exchange's Name\" must be specified")
	}

	if sess.Binding().RoutingKey == "" {
		return nil, fmt.Errorf("consumerBuild Routing error: the \"binding's RoutingKey\" must be specified")
	}

	err := sess.Establish(cb.rabbitMqConn())
	if err != nil {
		return nil, fmt.Errorf("consumerBuild Routing error: %w", err)
	}
	return routing.New(sess), nil
}

func (cb *consumerBuild) Topic() (external.Consumer, error) {
	cb.sessionOptions = append(cb.sessionOptions, session.WithBrokerOptions(
		broker.WithQueue(queue.SetName(""),
			queue.SetExclusive(true),
		),
		broker.WithExchange(
			exchange.SetDurable(true),
			exchange.SetType(external.XExchangeTopic),
		),
	))
	sess := cb.sess()
	if sess.Exchange().Name == "" {
		return nil, fmt.Errorf("consumerBuild Topic error: the \"exchange's Name\" must be specified")
	}
	if sess.Binding().RoutingKey == "" {
		return nil, fmt.Errorf("consumerBuild Topic error: the \"binding's RoutingKey\" must be specified")
	}

	err := sess.Establish(cb.rabbitMqConn())
	if err != nil {
		return nil, fmt.Errorf("consumerBuild Routing error: %w", err)
	}
	return topic.New(sess), nil
}
