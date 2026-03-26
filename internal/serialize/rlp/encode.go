package rlp

import (
	"bytes"
	"encoding/binary"
	"reflect"
)

func Encode(v interface{}) ([]byte, error) {
	return encode(reflect.ValueOf(v)), nil
}

func encode(v reflect.Value) []byte {
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return encodeBytes([]byte{})
		}
		return encode(v.Elem())
	}

	switch v.Kind() {
	case reflect.String:
		return encodeBytes([]byte(v.String()))

	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return encodeBytes(v.Bytes())
		}
		return encodeListFromSlice(v)

	case reflect.Array:
		return encodeListFromSlice(v)

	case reflect.Struct:
		return encodeStruct(v)

	case reflect.Uint, reflect.Uint64:
		return encodeBytes(intToBytes(v.Uint()))

	case reflect.Int, reflect.Int64:
		return encodeBytes(intToBytes(uint64(v.Int())))

	case reflect.Bool:
		if v.Bool() {
			return encodeBytes([]byte{1})
		}
		return encodeBytes([]byte{})

	default:
		return encodeBytes([]byte{})
	}
}

func encodeStruct(v reflect.Value) []byte {
	var buf bytes.Buffer

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)

		if field.PkgPath != "" {
			continue
		}

		buf.Write(encode(v.Field(i)))
	}

	return encodeList(buf.Bytes())
}

func encodeListFromSlice(v reflect.Value) []byte {
	var buf bytes.Buffer

	for i := 0; i < v.Len(); i++ {
		buf.Write(encode(v.Index(i)))
	}

	return encodeList(buf.Bytes())
}

func encodeBytes(b []byte) []byte {
	l := len(b)

	if l == 1 && b[0] < 0x80 {
		return b
	}

	if l <= 55 {
		return append([]byte{byte(0x80 + l)}, b...)
	}

	lenBytes := intToBytes(uint64(l))
	return append(append([]byte{byte(0xb7 + len(lenBytes))}, lenBytes...), b...)
}

func encodeList(b []byte) []byte {
	l := len(b)

	if l <= 55 {
		return append([]byte{byte(0xc0 + l)}, b...)
	}

	lenBytes := intToBytes(uint64(l))
	return append(append([]byte{byte(0xf7 + len(lenBytes))}, lenBytes...), b...)
}

func intToBytes(i uint64) []byte {
	if i == 0 {
		return []byte{}
	}

	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], i)

	i2 := 0
	for i2 < len(buf) && buf[i2] == 0 {
		i2++
	}
	return buf[i2:]
}
