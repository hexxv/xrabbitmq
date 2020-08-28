package json

import (
	"encoding/json"
)

// Codec implements the codec.Codec interface
type Codec struct{}

// NewCodec returns a new Serializer.
func NewCodec() *Codec {
	return &Codec{}
}

// Marshal returns the JSON encoding of v.
func (s *Codec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal parses the JSON-encoded data and stores the result
// in the value pointed to by v.
func (s *Codec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
