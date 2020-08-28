package xrabbitmq

type ConfOption func(*ConfOptions)

// RabbitMQ 建立链接所需配置项
type ConfOptions struct {
	// rabbitMQ地址
	Host string

	// rabbitMQ端口
	Port int

	// rabbitMQ用户名
	User string

	// rabbitMQ密码
	Pwd string

	// VHost 虚拟主机，一个broker里可以开设多个vhost，用作不用用户的权限分离
	VHost string
}

func defaultConfOptions(opts ...ConfOption) ConfOptions {
	opt := ConfOptions{
		Host:  "localhost", // 默认本机地址
		Port:  5672,        // RabbitMQ默认端口号
		User:  "guest",     // 默认用户名
		Pwd:   "guest",     // 默认用户密码
		VHost: "/",         // 默认虚拟主机地址
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func WithHost(host string) ConfOption {
	return func(options *ConfOptions) {
		options.Host = host
	}
}

func WithPort(port int) ConfOption {
	return func(options *ConfOptions) {
		options.Port = port
	}
}

func WithUser(user string) ConfOption {
	return func(options *ConfOptions) {
		options.User = user
	}
}

func WithPwd(pwd string) ConfOption {
	return func(options *ConfOptions) {
		options.Pwd = pwd
	}
}

func WithVHost(vh string) ConfOption {
	return func(options *ConfOptions) {
		options.VHost = vh
	}
}
