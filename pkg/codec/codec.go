package codec

import (
	"xrabbitmq/pkg/codec/json"
)

var defaultCodec Codec = json.NewCodec()

func Marshal(v interface{}) ([]byte, error) {
	return defaultCodec.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return defaultCodec.Unmarshal(data, v)
}

func SetCodec(c Codec) {
	defaultCodec = c
}
