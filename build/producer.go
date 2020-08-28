package build

import (
	"fmt"
	"xrabbitmq/pkg/external"
	"xrabbitmq/pkg/producer/publish"
	"xrabbitmq/pkg/producer/routing"
	"xrabbitmq/pkg/producer/simple"
	"xrabbitmq/pkg/producer/topic"
	"xrabbitmq/pkg/producer/work"
	"xrabbitmq/pkg/session"
	"xrabbitmq/pkg/session/broker"
	"xrabbitmq/pkg/session/broker/binding"
	"xrabbitmq/pkg/session/broker/exchange"
	"xrabbitmq/pkg/session/broker/queue"
)

type producerBuild struct {
	sessionOptions []session.Option
	buildRequired  Required
}

func NewProducerBuild(br Required, sos ...session.Option) *producerBuild {
	return &producerBuild{
		sessionOptions: sos,
		buildRequired:  br,
	}
}

func (cb *producerBuild) sess() *session.Session {
	return session.NewSession(cb.sessionOptions...)
}

func (cb *producerBuild) rabbitMqConn() *external.XConnection {
	builder := depend{}
	cb.buildRequired(&builder)
	return builder.conn
}

func (cb *producerBuild) Simple() (external.Producer, error) {
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
		return nil, fmt.Errorf("producerBuild Simple error: the \"queue's Name\" must be specified")
	}
	err := sess.Establish(cb.rabbitMqConn())
	if err != nil {
		return nil, fmt.Errorf("producerBuild Simple error: %w", err)
	}
	return simple.New(sess), nil
}

func (cb *producerBuild) Work() (external.Producer, error) {
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
		return nil, fmt.Errorf("producerBuild Work error: the \"queue's Name\" must be specified")
	}
	err := sess.Establish(cb.rabbitMqConn())
	if err != nil {
		return nil, fmt.Errorf("producerBuild Work error: %w", err)
	}
	return work.New(sess), nil
}

func (cb *producerBuild) Publish() (external.Producer, error) {
	cb.sessionOptions = append(cb.sessionOptions, session.WithBrokerOptions(
		broker.WithQueue(queue.SetName(""),
			queue.SetExclusive(true),
		),
		broker.WithExchange(
			exchange.SetDurable(true),
			exchange.SetType(external.XExchangeFanout),
		),
	))
	sess := cb.sess()

	if sess.Exchange().Name == "" {
		return nil, fmt.Errorf("producerBuild Publish error: the \"exchange's Name\" must be specified")
	}
	err := sess.Establish(cb.rabbitMqConn())
	if err != nil {
		return nil, fmt.Errorf("producerBuild Publish error: %w", err)
	}
	return publish.New(sess), nil
}

func (cb *producerBuild) Routing(dynamic bool) (external.Producer, error) {
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
		return nil, fmt.Errorf("producerBuild Routing error: the \"exchange's Name\" must be specified")
	}

	if dynamic {
		err := sess.Establish(cb.rabbitMqConn())
		if err != nil {
			return nil, fmt.Errorf("producerBuild Routing Dynamic error: %w", err)
		}
		return routing.NewDynamic(sess), nil
	} else {
		if sess.Binding().RoutingKey == "" {
			return nil, fmt.Errorf("producerBuild Routing error: the \"binding's RoutingKey\" must be specified")
		}
		err := sess.Establish(cb.rabbitMqConn())
		if err != nil {
			return nil, fmt.Errorf("producerBuild Routing error: %w", err)
		}
		return routing.New(sess), nil
	}
}

func (cb *producerBuild) Topic(dynamic bool) (external.Producer, error) {
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
		return nil, fmt.Errorf("producerBuild Topic error: the \"exchange's Name\" must be specified")
	}

	if dynamic {
		err := sess.Establish(cb.rabbitMqConn())
		if err != nil {
			return nil, fmt.Errorf("producerBuild Topic Dynamic error: %w", err)
		}
		return topic.NewDynamic(sess), nil
	} else {
		if sess.Binding().RoutingKey == "" {
			return nil, fmt.Errorf("producerBuild Topic error: the \"binding's RoutingKey\" must be specified")
		}
		err := sess.Establish(cb.rabbitMqConn())
		if err != nil {
			return nil, fmt.Errorf("producerBuild Topic error: %w", err)
		}
		return topic.New(sess), nil
	}
}
