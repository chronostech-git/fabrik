package serialize

import (
	"bytes"
	"reflect"
)

type Serializable interface {
	Encode() ([]byte, error)
	Decode([]byte) error
}

func Encode(v any) ([]byte, error) {
	if s, ok := v.(Serializable); ok {
		return s.Encode()
	}

	buf := new(bytes.Buffer)
	if err := encodeValue(buf, reflect.ValueOf(v)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decode(data []byte, v any) error {
	if s, ok := v.(Serializable); ok {
		return s.Decode(data)
	}

	buf := bytes.NewReader(data)
	return decodeValue(buf, reflect.ValueOf(v).Elem())
}
