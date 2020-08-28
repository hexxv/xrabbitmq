package protobuf

import (
	"errors"
	"github.com/golang/protobuf/proto"
)

// ErrWrongValueType is the error used for marshal the value with protobuf encoding.
var ErrWrongValueType = errors.New("protobuf: convert on wrong type value")

// Codec implements the codec.Codec interface
type Codec struct{}

// Codec returns a new Codec.
func NewCodec() *Codec {
	return &Codec{}
}

// Marshal returns the protobuf encoding of v.
func (s *Codec) Marshal(v interface{}) ([]byte, error) {
	pb, ok := v.(proto.Message)
	if !ok {
		return nil, ErrWrongValueType
	}
	return proto.Marshal(pb)
}

// Unmarshal parses the protobuf-encoded data and stores the result
// in the value pointed to by v.
func (s *Codec) Unmarshal(data []byte, v interface{}) error {
	pb, ok := v.(proto.Message)
	if !ok {
		return ErrWrongValueType
	}
	return proto.Unmarshal(data, pb)
}
