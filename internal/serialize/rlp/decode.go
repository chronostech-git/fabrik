package rlp

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/chronostech-git/fabrik/internal/types"
)

func Decode(data []byte, out any) error {
	v := reflect.ValueOf(out)

	if v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.New("output must be pointer")
	}

	_, err := decodeValue(data, v.Elem())
	return err
}

func decodeValue(data []byte, v reflect.Value) ([]byte, error) {

	if len(data) == 0 {

		switch v.Kind() {

		case reflect.Slice:
			if v.Type().Elem().Kind() == reflect.Uint8 {
				v.SetBytes(nil)
				return data, nil
			}

		case reflect.Array:
			if v.Type().Elem().Kind() == reflect.Uint8 {
				for i := 0; i < v.Len(); i++ {
					v.Index(i).SetUint(0)
				}
				return data, nil
			}

		case reflect.String:
			v.SetString("")
			return data, nil
		}

		return nil, fmt.Errorf("rlp: empty input")
	}

	if v.Kind() == reflect.Ptr {

		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}

		return decodeValue(data, v.Elem())
	}

	if v.Type().Name() == "Amount" {

		var b []byte

		data, err := decodeValue(data, reflect.ValueOf(&b).Elem())
		if err != nil {
			return nil, err
		}

		v.Set(reflect.ValueOf(types.BytesToAmount(b)))

		return data, nil
	}

	prefix := data[0]

	switch {

	case prefix <= 0x7f:
		return setBytes(v, []byte{prefix}, data[1:])

	case prefix <= 0xb7:
		l := int(prefix - 0x80)
		b := data[1 : 1+l]
		return setBytes(v, b, data[1+l:])

	case prefix <= 0xbf:
		lenLen := int(prefix - 0xb7)
		l := bytesToInt(data[1 : 1+lenLen])
		start := 1 + lenLen
		b := data[start : start+l]
		return setBytes(v, b, data[start+l:])

	case prefix <= 0xf7:
		l := int(prefix - 0xc0)
		return decodeList(data[1:1+l], data[1+l:], v)

	default:
		lenLen := int(prefix - 0xf7)
		l := bytesToInt(data[1 : 1+lenLen])
		start := 1 + lenLen
		return decodeList(data[start:start+l], data[start+l:], v)
	}
}

func decodeList(listData []byte, rest []byte, v reflect.Value) ([]byte, error) {

	switch v.Kind() {

	case reflect.Struct:

		for i := 0; i < v.NumField(); i++ {

			if !v.Type().Field(i).IsExported() {
				continue
			}

			var err error

			listData, err = decodeValue(listData, v.Field(i))
			if err != nil {
				return nil, err
			}
		}

		return rest, nil

	case reflect.Slice:

		elemType := v.Type().Elem()

		slice := reflect.MakeSlice(v.Type(), 0, 0)

		for len(listData) > 0 {

			elem := reflect.New(elemType).Elem()

			before := len(listData)

			var err error

			listData, err = decodeValue(listData, elem)
			if err != nil {
				return nil, err
			}

			slice = reflect.Append(slice, elem)

			if len(listData) == before {
				return nil, fmt.Errorf("rlp: no progress in decode")
			}
		}

		v.Set(slice)

		return rest, nil

	default:
		return nil, fmt.Errorf("unsupported type in decodeList: %s", v.Kind())
	}
}

func setBytes(v reflect.Value, b []byte, rest []byte) ([]byte, error) {

	switch v.Kind() {

	case reflect.Bool:

		if len(b) == 0 {
			v.SetBool(false)
		} else if len(b) == 1 && b[0] == 1 {
			v.SetBool(true)
		} else {
			return nil, fmt.Errorf("invalid RLP bool encoding")
		}

		return rest, nil

	case reflect.Array:

		if v.Type().Elem().Kind() == reflect.Uint8 {

			if len(b) != v.Len() {
				return nil, fmt.Errorf("byte array length mismatch: got %d, want %d", len(b), v.Len())
			}

			for i := 0; i < v.Len(); i++ {
				v.Index(i).SetUint(uint64(b[i]))
			}

			return rest, nil
		}

	case reflect.Slice:

		if v.Type().Elem().Kind() == reflect.Uint8 {
			v.SetBytes(b)
			return rest, nil
		}

	case reflect.String:

		v.SetString(string(b))
		return rest, nil

	case reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uint:

		v.SetUint(uint64(bytesToInt(b)))
		return rest, nil

	case reflect.Int, reflect.Int64, reflect.Int32:

		v.SetInt(int64(bytesToInt(b)))
		return rest, nil
	}

	return nil, fmt.Errorf("unsupported type in setBytes: %s", v.Kind())
}

func bytesToInt(b []byte) int {

	n := 0

	for _, v := range b {
		n = (n << 8) | int(v)
	}

	return n
}
