package exchange

import (
	"xrabbitmq/pkg/external"
)

type Option func(exchange *Exchange)

// 交换机: 接收消息, 并根据路由键转发消息所绑定的队列
type Exchange struct {
	// Name: 交换机名字
	// 如果设置为空，则默认发送到RabbitMQ 默认的交换器中
	Name string

	// Typ: 交换机类型  direct/topic/fanout/headers
	//
	// 1. Direct Exchange
	// <1> 所有发送到Direct Exchange的消息被转发到RoutingKey中指定的Queue
	// 注意 : Direct模式可以使用RabbitMQ自带的Exchange(default Exchange),
	// 所以不需要将Exchange进行任何绑定(binding)操作, 消息传递时, RoutingKey必须完全匹配才会
	// 被队列接收, 否则该消息会被抛弃
	//
	// 2. Topic Exchange
	// <1> 所有发送到Topic Exchange的消息将被转发到所有关心RoutingKey中指定Topic的Queue上
	// <2> Exchange将RoutingKey和某个Topic进行模糊匹配, 此时队列需要绑定一个Topic
	//
	// 3. Fanout Exchange
	// <1> 不处理路由键, 只需要简单的将队列绑定到交换机上
	// <2> 发送到交换机的消息都会被转发到与该交换机绑定的所有队列上
	// <3> Fanout交换机转发消息是最快的
	//
	// 4. Headers Exchange
	// <1> Headers Exchange不使用RoutingKey去绑定, 而是通过消息headers的键值对匹配
	// <2> 这个Exchange很少会使用
	Typ string

	// Durable: 是否持久化
	Durable bool

	// AutoDelete: 是否自动删除：当最后一个绑定到Exchange上的队列删除后, 自动删除该Exchange
	AutoDelete bool

	// Internal: 交换机声明为“内部”不接受接受发布。
	// 当前Exchange是否用于RabbitMQ内部使用, 默认为False, 这个属性很少会用到
	// 当希望实现不应该向代理的用户公开的相互交换拓扑时，内部交换非常有用。
	Internal bool

	// NoWait: 是否阻塞处理
	NoWait bool

	// Args: 额外的属性
	// 对于需要额外参数的交换类型，可以发送特定于服务器的交换实现的参数表
	// 扩展参数, 用于扩展AMQP协议制定化使用
	Args external.XTable
}

func SetName(name string) Option {
	return func(exchange *Exchange) {
		exchange.Name = name
	}
}

func SetType(typ string) Option {
	return func(exchange *Exchange) {
		exchange.Typ = typ
	}
}

func SetDurable(durable bool) Option {
	return func(exchange *Exchange) {
		exchange.Durable = durable
	}
}

func SetAutoDelete(autoDel bool) Option {
	return func(exchange *Exchange) {
		exchange.AutoDelete = autoDel
	}
}

func SetInternal(internal bool) Option {
	return func(exchange *Exchange) {
		exchange.Internal = internal
	}
}

func SetNoWait(noWait bool) Option {
	return func(exchange *Exchange) {
		exchange.NoWait = noWait
	}
}

func WithArgs(args external.XTable) Option {
	return func(exchange *Exchange) {
		if exchange.Args == nil {
			exchange.Args = make(external.XTable)
		}
		for k, v := range args {
			exchange.Args[k] = v
		}
	}
}
