package utils

import (
	"fmt"
	"xrabbitmq/pkg/external"
)

// 关闭RabbitMQ Channel
func CancelChannel(channel *external.XChannel, tag string) error {
	if err := channel.Cancel(tag, true); err != nil {
		if amqpError, ok := err.(*external.XError); ok {
			if amqpError.Code != external.ChannelError {
				return fmt.Errorf("AMQP connection channel close error: %w", err)
			}
		}
	}
	if err := channel.Close(); err != nil {
		return err
	}
	return nil
}
