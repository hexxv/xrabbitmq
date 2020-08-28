package consumeropts

import (
	"xrabbitmq/pkg/external"
)

type Option func(*Options)

// 消费者配置项
type Options struct {
	// Tag: 消费者由一个字符串标识，该字符串是唯一的，并且适用于此通道上的所有消费者。
	Tag string

	// AutoAck: 是否自动确认
	AutoAck bool

	// Exclusive: 是否具有排他性
	Exclusive bool

	// NoLocal: 当noLocal为true时，服务器将不会将从同一连接发送的发布传递给此使用者。(不要从同一通道使用发布和消费)
	NoLocal bool

	// NoWait: 是否阻塞处理
	NoWait bool

	// Args: 额外的属性
	Args external.XTable
}

func SetTag(tag string) Option {
	return func(options *Options) {
		options.Tag = tag
	}
}

func SetAutoAck(autoAck bool) Option {
	return func(options *Options) {
		options.AutoAck = autoAck
	}
}

func SetExclusive(exclusive bool) Option {
	return func(options *Options) {
		options.Exclusive = exclusive
	}
}

func SetNoLocal(noLocal bool) Option {
	return func(options *Options) {
		options.NoLocal = noLocal
	}
}

func SetNoWait(noWait bool) Option {
	return func(options *Options) {
		options.NoWait = noWait
	}
}

func WithArgs(args external.XTable) Option {
	return func(options *Options) {
		if options.Args == nil {
			options.Args = make(external.XTable)
		}
		for k, v := range args {
			options.Args[k] = v
		}
	}
}
