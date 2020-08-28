package binding

import (
	"xrabbitmq/pkg/external"
)

type Option func(*Binding)

// 绑定选项
type Binding struct {
	// RoutingKey: 使用匹配的路由键将消息发布到给定队列——每个队列都有一个默认绑定到默认交换器，
	// 使用它们的队列名称，这样就可以通过默认交换器将消息发送到队列
	RoutingKey string

	// NoWait: 是否阻塞处理
	NoWait bool

	// Args: 额外的属性
	Args external.XTable
}

func SetRoutingKey(key string) Option {
	return func(options *Binding) {
		options.RoutingKey = key
	}
}

func SetNoWait(noWait bool) Option {
	return func(options *Binding) {
		options.NoWait = noWait
	}
}

func WithArgs(args external.XTable) Option {
	return func(options *Binding) {
		if options.Args == nil {
			options.Args = make(external.XTable)
		}
		for k, v := range args {
			options.Args[k] = v
		}
	}
}
