package produceropts

type Option = func(*Options)

type Options struct {
	// RoutingKey: 使用匹配的路由键将消息发布到给定队列——每个队列都有一个默认绑定到默认交换器，
	// 使用它们的队列名称，这样就可以通过默认交换器将消息发送到队列
	RoutingKey string

	// Tag:  tag
	Tag string

	// 概括来说:
	//	1. mandatory标志告诉服务器至少将该消息route到一个队列中，否则将消息返还给生产者
	//
	//	2. immediate标志告诉服务器如果该消息关联的queue上有消费者，则马上将消息投递给它，
	//	   如果所有queue都没有消费者，直接把消息返还给生产者，不用将消息入队列等待消费者了。
	//
	// Mandatory: Queue should be on the server/broker
	//
	// 当mandatory标志位设置为true时，如果exchange根据自身类型和消息routeKey无法找到一个符
	// 合条件的queue，那么会调用basic.return方法将消息返还给生产者；当mandatory设为false时，
	// 出现上述情形broker会直接将消息扔掉。
	Mandatory bool

	// Immediate: Consumer should be bound to server
	//
	// 当immediate标志位设置为true时，如果exchange在将消息route到queue(s)时发现对应的queue
	// 上没有消费者，那么这条消息不会放入队列中。当与消息routeKey关联的所有queue(一个或多个)都
	// 没有消费者时，该消息会通过basic.return方法返还给生产者。
	//
	// 注意：在RabbitMQ3.0以后的版本里，去掉了immediate参数支持，发送带immediate标记的publish会返回如下错误：
	// “{amqp_error,not_implemented,"immediate=true",'basic.publish'}”
	// 对此官方给出的解释是:
	// Removal of "immediate" flag
	//
	// What changed? We removed support for the rarely-used "immediate" flag on AMQP's basic.publish.
	// Why on earth did you do that? Support for "immediate" made many parts of the codebase more complex,
	// particularly around mirrored queues. It also stood in the way of our being able to deliver
	// substantial performance improvements in mirrored queues.
	//
	// What do I need to do? If you just want to be able to publish messages that will be dropped if
	// they are not consumed immediately, you can publish to a queue with a TTL of 0.
	// If you also need your publisher to be able to determine that this has happened, you can also use the
	// DLX feature to route such messages to another queue, from which the publisher can consume them.
	//
	// 大概意思就是: immediate标记会影响镜像队列性能，增加代码复杂性，并建议采用"设置消息TTL"和"DLX"等方式替代。
	// Immediate bool
}

func SetRoutingKey(key string) Option {
	return func(options *Options) {
		options.RoutingKey = key
	}
}

func SetTag(tag string) Option {
	return func(options *Options) {
		options.Tag = tag
	}
}

func SetMandatory(mandatory bool) Option {
	return func(options *Options) {
		options.Mandatory = mandatory
	}
}

// 注意：在RabbitMQ3.0以后的版本里，去掉了immediate参数支持，发送带immediate标记的publish会返回如下错误：
// “{amqp_error,not_implemented,"immediate=true",'basic.publish'}”
// 对此官方给出的解释是:
// Removal of "immediate" flag
//
// What changed? We removed support for the rarely-used "immediate" flag on AMQP's basic.publish.
// Why on earth did you do that? Support for "immediate" made many parts of the codebase more complex,
// particularly around mirrored queues. It also stood in the way of our being able to deliver
// substantial performance improvements in mirrored queues.
//
// What do I need to do? If you just want to be able to publish messages that will be dropped if
// they are not consumed immediately, you can publish to a queue with a TTL of 0.
// If you also need your publisher to be able to determine that this has happened, you can also use the
// DLX feature to route such messages to another queue, from which the publisher can consume them.
//
// 大概意思就是: immediate标记会影响镜像队列性能，增加代码复杂性，并建议采用"设置消息TTL"和"DLX"等方式替代。
// func SetImmediate(immediate bool) Option {
// 	return func(options *Options) {
// 		options.Immediate = immediate
// 	}
// }
