package queue

import (
	"xrabbitmq/pkg/external"
)

type Option func(queue *Queue)

// 消息队列, 实际存储消息数据
type Queue struct {
	// Name: 队列名字，可能为空，此时服务器将生成唯一的名称，该名称将在队列结构的名称字段中返回。
	Name string

	// Durable: 是否持久化
	// 为true 则设置队列为持久化。持久化的队列会存盘，在服务器重启的时候可以保证不丢失相关信息
	Durable bool

	// AutoDelete: 是否自动删除
	// 当最后一个监听被移除后，自动删除队列；也就是说至少有一个消费者连接到这个队列，之后所有与这个
	// 队列连接的消费者都断开时，才会自动删除
	AutoDelete bool

	// Exclusive: 是否具有排他性
	// 为true 则设置队列为排他的。如果一个队列被声明为排他队列，该队列仅对首次声明它的连接可见，
	// 并在连接断开时自动删除。
	// 这里需要注意三点:
	// 1. 排他队列是基于连接(Connection) 可见的，同一个连接的不同信道(Channel)是可以同时访问
	//    同一连接创建的排他队列;
	// 2. "首次"是指如果一个连接己经声明了一个排他队列，其他连接是不允许建立同名的排他队列的，这个
	//    与普通队列不同:
	// 3. 即使该队列是持久化的，一旦连接关闭或者客户端退出，该排他队列都会被自动删除，这种队列适用于
	//    一个客户端同时发送和读取消息的应用场景
	Exclusive bool

	// NoWait: 是否阻塞处理
	// 当noWait为true时，队列将假定在服务器上声明。如果现有队列的条件满足，或者试
	// 图从不同的连接修改现有队列，则会出现通道异常。
	NoWait bool

	// Args: 额外的属性
	Args external.XTable
}

func SetName(name string) Option {
	return func(queue *Queue) {
		queue.Name = name
	}
}

func SetDurable(durable bool) Option {
	return func(queue *Queue) {
		queue.Durable = durable
	}
}

func SetAutoDelete(autoDel bool) Option {
	return func(queue *Queue) {
		queue.AutoDelete = autoDel
	}
}

func SetExclusive(exclusive bool) Option {
	return func(queue *Queue) {
		queue.Exclusive = exclusive
	}
}

func SetNoWait(noWait bool) Option {
	return func(queue *Queue) {
		queue.NoWait = noWait
	}
}

func WithArgs(args external.XTable) Option {
	return func(queue *Queue) {
		if queue.Args == nil {
			queue.Args = make(external.XTable)
		}
		for k, v := range args {
			queue.Args[k] = v
		}
	}
}
