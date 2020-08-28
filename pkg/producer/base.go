package producer

import (
	"fmt"
	"github.com/streadway/amqp"
	"xrabbitmq/pkg/external"
	"xrabbitmq/pkg/internal/utils"
	"xrabbitmq/pkg/log"
	"xrabbitmq/pkg/session"
	"sync"
	"time"
)

type Producer struct {
	// session 消费者与交换机/队列/binding的声明保障
	// 消费者会占用一个RabbitMQ的通信管道，不要忘记使用 Consumer.Cancel() 来释放它
	session *session.Session

	// cancelOnce 关闭通信管道幂等
	cancelOnce sync.Once

	// deliveries all deliveries from server will send to this channel
	messages <-chan *external.XPublishMsg

	// done: a notifiyng channel for publishings
	// will be used for sync. between close channel and consume handler
	done chan error

	// producer model simple/work/publish/routing/topic
	model Model
}

func NewProducer(sess *session.Session, mod Model) *Producer {
	return &Producer{
		session: sess,
		done:    make(chan error),
		model:   mod,
	}
}

func (p *Producer) Sess() *session.Session {
	return p.session
}

func (p *Producer) Cancel() error {
	return p.releaseChannel()
}

// Producer 在处理完所有消息后 断掉连接
func (p *Producer) releaseChannel() (err error) {
	defer func() {
		p.cancelOnce.Do(func() {
			err = utils.CancelChannel(p.session.Channel(), p.session.OptionsProducer().Tag)
		})
	}()

	if p.messages == nil {
		close(p.done)
	}

	return <-p.done
}

func (p *Producer) NotifyReturn(handleFunc external.ReturnHandleFunc) {
	go func() {
		defer func() {
			if x := recover(); x != nil {
				log.Logger.Errorf("producer listen panic: %+v", x)
			}
		}()
		for msg := range p.session.Channel().NotifyReturn(make(chan external.XReturn, 1)) {
			handleFunc(msg)
		}
	}()
}

// 监听生产者的消息确认
// func (p *Producer) ListenNotify(v interface{}) error {
// 	// 监听消息未到达queue的回调
// 	returnNotify, returnHandler := p.notifyReturn(v)
//
// 	//
// 	publishNotify, publishHandler := p.notifyPublish(v)
// 	if returnHandler != nil || publishHandler != nil {
// 		p.wg.Add(1)
// 		go func() {
// 			defer func() {
// 				if x := recover(); x != nil {
// 					log.Logger.Errorf("producer listen panic: %+v", x)
// 				}
// 				p.wg.Done()
// 			}()
// 			for {
// 				select {
// 				case msg := <-returnNotify:
// 					returnHandler(msg)
// 				case msg := <-publishNotify:
// 					publishHandler(msg)
// 				case <-p.canceling:
// 					log.Logger.Warning("producer shutdown, listener return now.")
// 					return
// 				}
// 			}
// 		}()
// 		return nil
// 	} else {
// 		return fmt.Errorf("producer listen 希望传入一个结构体，它具有" +
// 			"\"NotifyReturn(message external.XReturn)\" 或者 \"NotifyPublish(message external.XConfirmation)\"的方法" +
// 			"或者直接传入一个像: \"func(message external.XReturn)\" 或者 \"func(message external.XConfirmation)\"这样的函数" +
// 			"但是传入的interface没有满足任何一个条件，所以\"ListenNotify\"将不会运行。")
// 	}
// }
//
// // NotifyReturn 当发消息时交换机不在，或者队列不在，或者Routing-key错误， 那么相应的Renting-Listenner就可以工作了，
// // 注意，使用这个功能时，要把basicpublish中使用 mandatory开启
// func (p *Producer) notifyReturn(v interface{}) (chan external.XReturn, func(message external.XReturn)) {
// 	ch := make(chan external.XReturn)
// 	switch handler := v.(type) {
// 	case external.ReturnHandleFunc:
// 		return p.session.Channel().NotifyReturn(ch), handler
// 	case external.ProducerHandler:
// 		return p.session.Channel().NotifyReturn(ch), handler.NotifyReturn
// 	case external.ReturnHandler:
// 		return p.session.Channel().NotifyReturn(ch), handler.NotifyReturn
// 	default:
// 		return ch, nil
// 	}
// }
//
// func (p *Producer) notifyPublish(v interface{}) (chan external.XConfirmation, func(message external.XConfirmation)) {
// 	ch := make(chan external.XConfirmation)
// 	switch handler := v.(type) {
// 	case external.PublishHandleFunc:
// 		return p.session.Channel().NotifyPublish(ch), handler
// 	case external.ProducerHandler:
// 		return p.session.Channel().NotifyPublish(ch), handler.NotifyPublish
// 	case external.PublishHandler:
// 		return p.session.Channel().NotifyPublish(ch), handler.NotifyPublish
// 	default:
// 		return ch, nil
// 	}
// }

// var msg amqp.Publishing

// amqp.Publishing {
//        // Application or exchange specific fields,
//        // the headers exchange will inspect this field.
//        Headers Table

//        // Properties
//        ContentType     string    // MIME content type
//        ContentEncoding string    // MIME content encoding
//        DeliveryMode    uint8     // Transient (0 or 1) or Persistent (2)
//        Priority        uint8     // 0 to 9
//        CorrelationId   string    // correlation identifier
//        ReplyTo         string    // address to to reply to (ex: RPC)
//        Expiration      string    // message expiration spec
//        MessageId       string    // message identifier
//        Timestamp       time.Time // message timestamp
//        Type            string    // message type name
//        UserId          string    // creating user id - ex: "guest"
//        AppId           string    // creating application id

//        // The application specific payload of the message
//        Body []byte
// }

// msg.Body, err = codec.Marshal(v)
// if err != nil {
// 	log.Logger.Errorf("simple.producer.Marshal error: %s", err)
// 	return err
// }

func (p *Producer) Model() Model {
	return p.model
}

func (p *Producer) Done(err error) {
	p.done <- err
}

func (p *Producer) key(customKey string) string {
	switch p.model {
	case ModelTopic, ModelRouting:
		if customKey != "" {
			return customKey
		} else {
			return p.session.Binding().RoutingKey
		}
	case ModelRoutingDynamic, ModelTopicDynamic:
		return customKey
	case ModelWork, ModelSimple:
		return p.session.Queue().Name
	case ModelPublish:
		return ""
	}
	log.Logger.Error("producer get key error: %s", p.model)
	return ""
}

func (p *Producer) Publish(messages <-chan *external.XPublishMsg) {
	p.messages = messages
	var (
		exchangeOptions = p.session.Exchange()
		producerOptions = p.session.OptionsProducer()
		pub             = p.session.Channel()
		reading         = p.messages
		pending         = make(chan *external.XPublishMsg, 1)
		confirm         = make(chan external.XConfirmation, 1)
	)

	if err := pub.Confirm(false); err != nil {
		log.Logger.Error("publisher confirms not supported")
		close(confirm)
	} else {
		pub.NotifyPublish(confirm)
	}

	log.Logger.Info("publishing...")

Publish:
	for {
		var body *external.XPublishMsg
		select {
		case confirmed, ok := <-confirm:
			if !ok {
				break Publish
			}
			if !confirmed.Ack {
				log.Logger.Error("nack message %d, body: %q", confirmed.DeliveryTag, string(body.Body))
			}
			reading = messages

		case body = <-pending:
			if body == nil {
				return
			}
			fmt.Println("recv body:", string(body.Body))

			// 模拟复杂度
			time.Sleep(time.Millisecond * 20)

			err := pub.Publish(
				exchangeOptions.Name,
				p.key(body.RoutingKey),
				producerOptions.Mandatory,
				false,
				amqp.Publishing{
					ContentType: body.ContentType,
					Body:        body.Body,
				})
			// Retry failed delivery
			if err != nil {
				log.Logger.Errorf("producer send error: %s will retry after 1s", err)
				time.Sleep(time.Second * 1)
				pending <- body
				break Publish
			}

		case body = <-reading:
			// fmt.Println("recv pub body:", string(body))
			// all messages consumed
			if body == nil {
				fmt.Println("recv nil, return now.")
				return
			}
			// work on pending delivery until ack'd
			pending <- body
			reading = nil
		}
	}
}
