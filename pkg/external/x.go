package external

import "github.com/streadway/amqp"

// 等价替换：目的是为了让外部包/文件在使用xrabbitmq的时候不用导入"github.com/streadway/amqp"

type (
	XConnection = amqp.Connection

	XPublishing = amqp.Publishing

	XReturn = amqp.Return

	XDelivery = amqp.Delivery

	XChannel = amqp.Channel

	XError = amqp.Error

	XURI = amqp.URI

	XBlocking = amqp.Blocking

	XTable = amqp.Table

	XConfirmation = amqp.Confirmation

	XQueue = amqp.Queue
)

type XPublishMsg struct {
	RoutingKey  string // 路由key
	ContentType string // "text/plain" "application/octet-stream" 可以不填
	Body        []byte // data
}

func XDial(url string) (*XConnection, error) {
	return amqp.Dial(url)
}

const (
	ContentTooLarge    = amqp.ContentTooLarge
	NoRoute            = amqp.NoRoute
	NoConsumers        = amqp.NoConsumers
	ConnectionForced   = amqp.ConnectionForced
	InvalidPath        = amqp.InvalidPath
	AccessRefused      = amqp.AccessRefused
	NotFound           = amqp.NotFound
	ResourceLocked     = amqp.ResourceLocked
	PreconditionFailed = amqp.PreconditionFailed
	FrameError         = amqp.FrameError
	SyntaxError        = amqp.SyntaxError
	CommandInvalid     = amqp.CommandInvalid
	ChannelError       = amqp.ChannelError
	UnexpectedFrame    = amqp.UnexpectedFrame
	ResourceError      = amqp.ResourceError
	NotAllowed         = amqp.NotAllowed
	NotImplemented     = amqp.NotImplemented
	InternalError      = amqp.InternalError
)

const (
	XExchangeDirect  = amqp.ExchangeDirect
	XExchangeFanout  = amqp.ExchangeFanout
	XExchangeTopic   = amqp.ExchangeTopic
	XExchangeHeaders = amqp.ExchangeHeaders
)
