package xrabbitmq

import (
	"github.com/sirupsen/logrus"
	"xrabbitmq/pkg/codec"
	"xrabbitmq/pkg/log"
)

// SetLogger 替换xrabbitmq包默认的logger
func SetLogger(iLogger logrus.FieldLogger) {
	log.SetLogger(iLogger)
}

// SetCodec 替换xrabbitmq包默认的解码器
func SetCodec(iCodec codec.Codec) {
	codec.SetCodec(iCodec)
}
