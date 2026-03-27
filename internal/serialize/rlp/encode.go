package rlp

import (
	"fmt"
	"reflect"

	"github.com/chronostech-git/fabrik/internal/types"
)

func Encode(input any) ([]byte, error) {
	v := reflect.ValueOf(input)
	return encodeValue(v)
}

func encodeValue(v reflect.Value) ([]byte, error) {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return []byte{0x80}, nil
		}
		return encodeValue(v.Elem())
	}

	if v.Type().Name() == "Amount" {
		amt := v.Interface().(types.Amount)
		return encodeBytes(amt.Bytes()), nil
	}

	switch v.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		b := intToBytes(int(v.Uint()))
		return encodeBytes(b), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		b := intToBytes(int(v.Int()))
		return encodeBytes(b), nil

	case reflect.String:
		return encodeBytes([]byte(v.String())), nil

	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return encodeBytes(v.Bytes()), nil
		}
		var encodedItems []byte
		for i := 0; i < v.Len(); i++ {
			e, err := encodeValue(v.Index(i))
			if err != nil {
				return nil, err
			}
			encodedItems = append(encodedItems, e...)
		}
		return encodeList(encodedItems), nil

	case reflect.Array: // <--- NEW: handle fixed-size arrays like [32]byte
		if v.Type().Elem().Kind() == reflect.Uint8 {
			b := make([]byte, v.Len())
			for i := 0; i < v.Len(); i++ {
				b[i] = byte(v.Index(i).Uint())
			}
			return encodeBytes(b), nil
		}
		var encodedItems []byte
		for i := 0; i < v.Len(); i++ {
			e, err := encodeValue(v.Index(i))
			if err != nil {
				return nil, err
			}
			encodedItems = append(encodedItems, e...)
		}
		return encodeList(encodedItems), nil

	case reflect.Struct:
		var encodedFields []byte
		for i := 0; i < v.NumField(); i++ {
			if !v.Type().Field(i).IsExported() {
				continue
			}
			e, err := encodeValue(v.Field(i))
			if err != nil {
				return nil, err
			}
			encodedFields = append(encodedFields, e...)
		}
		return encodeList(encodedFields), nil

	default:
		return nil, fmt.Errorf("unsupported type in encodeValue: %s", v.Kind())
	}
}

func encodeBytes(b []byte) []byte {
	if len(b) == 1 && b[0] <= 0x7f {
		return b
	}
	if len(b) <= 55 {
		return append([]byte{0x80 + byte(len(b))}, b...)
	}
	lenBytes := intToBytes(len(b))
	return append(append([]byte{0xb7 + byte(len(lenBytes))}, lenBytes...), b...)
}

func encodeList(b []byte) []byte {
	if len(b) <= 55 {
		return append([]byte{0xc0 + byte(len(b))}, b...)
	}
	lenBytes := intToBytes(len(b))
	return append(append([]byte{0xf7 + byte(len(lenBytes))}, lenBytes...), b...)
}

func intToBytes(n int) []byte {
	if n == 0 {
		return []byte{}
	}
	var b []byte
	for n > 0 {
		b = append([]byte{byte(n & 0xff)}, b...)
		n >>= 8
	}
	return b
}
