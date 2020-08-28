package codec

// 序列化
type Marshaler interface {
	Marshal(interface{}) ([]byte, error)
}

// 反序列化
type Unmarshaler interface {
	Unmarshal([]byte, interface{}) error
}

// 解码器
type Codec interface {
	Marshaler
	Unmarshaler
}
