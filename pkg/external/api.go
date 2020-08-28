// Expose the interface for external band calls
package external

// 生产者
type Producer interface {
	// NotifyReturn
	// exchange到queue成功,则不回调return
	// exchange到queue失败,则回调return(需设置mandatory=true,否则不回回调,消息就丢了)
	NotifyReturn(handleFunc ReturnHandleFunc)

	// Publish 生产一条消息
	Publish(messages <-chan *XPublishMsg) error

	// Cancel 关闭通信管道，释放资源
	Cancel() error
}

// 消费者
type Consumer interface {
	// Qos 即服务质量保证
	// 在非自动确认消息的前提下，如果一定数目的消息（通过基于consume或者
	// channel设置Qos的值）未被确认前，不进行消费新的消息。
	Qos(int) error

	// Consume 开始消费，阻塞式
	Consume(func(delivery XDelivery)) error

	// Cancel 关闭通信管道，释放资源
	Cancel() error
}

// 生产者的构建者
type ProducerBuilder interface {
	// 简单模式
	Simple() (Producer, error)

	// 工作模式
	Work() (Producer, error)

	// 订阅模式
	Publish() (Producer, error)

	// 路由模式
	Routing(dynamic bool) (Producer, error)

	// 话题模式
	Topic(dynamic bool) (Producer, error)
}

// 消费者的构建者
type ConsumerBuilder interface {
	// 简单模式
	Simple() (Consumer, error)

	// 工作模式
	Work() (Consumer, error)

	// 订阅模式
	Publish() (Consumer, error)

	// 路由模式
	Routing() (Consumer, error)

	// 话题模式
	Topic() (Consumer, error)
}

// exchange到queue成功,则不回调return
// exchange到queue失败,则回调return(需设置mandatory=true,否则不回回调,消息就丢了)
// 如果消息没有到exchange,则confirm回调,ack=false
// 如果消息到达exchange,则confirm回调,ack=true
type ReturnHandleFunc func(message XReturn)
