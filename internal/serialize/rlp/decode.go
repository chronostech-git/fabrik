package rlp

import (
	"errors"
	"math/big"
	"reflect"
)

func Decode(data []byte, out interface{}) error {
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Pointer {
		return errors.New("output must be pointer")
	}

	_, err := decode(data, v.Elem())
	return err
}

func decode(data []byte, v reflect.Value) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("empty input")
	}

	prefix := data[0]

	switch {
	case prefix <= 0x7f:
		return decodeBytes(data[:1], v, data[1:])

	case prefix <= 0xb7:
		l := int(prefix - 0x80)
		start := 1
		end := start + l
		return decodeBytes(data[start:end], v, data[end:])

	case prefix <= 0xbf:
		lenLen := int(prefix - 0xb7)
		l := bytesToInt(data[1 : 1+lenLen])
		start := 1 + lenLen
		end := start + int(l)
		return decodeBytes(data[start:end], v, data[end:])

	case prefix <= 0xf7:
		l := int(prefix - 0xc0)
		start := 1
		end := start + l
		return decodeList(data[start:end], v, data[end:])

	default:
		lenLen := int(prefix - 0xf7)
		l := bytesToInt(data[1 : 1+lenLen])
		start := 1 + lenLen
		end := start + int(l)
		return decodeList(data[start:end], v, data[end:])
	}
}

func decodeBytes(b []byte, v reflect.Value, rest []byte) ([]byte, error) {
	switch v.Kind() {
	case reflect.String:
		v.SetString(string(b))

	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			v.SetBytes(b)
		} else {
			return nil, errors.New("invalid slice type")
		}

	case reflect.Uint, reflect.Uint64:
		v.SetUint(bytesToInt(b))

	case reflect.Int, reflect.Int64:
		v.SetInt(int64(bytesToInt(b)))

	case reflect.Bool:
		v.SetBool(len(b) > 0)

	default:
		return nil, errors.New("unsupported type in decodeBytes")
	}

	return rest, nil
}

func decodeList(b []byte, v reflect.Value, rest []byte) ([]byte, error) {
	switch v.Kind() {
	case reflect.Slice:
		elemType := v.Type().Elem()
		slice := reflect.MakeSlice(v.Type(), 0, 0)

		for len(b) > 0 {
			elem := reflect.New(elemType).Elem()

			var err error
			b, err = decode(b, elem)
			if err != nil {
				return nil, err
			}

			slice = reflect.Append(slice, elem)
		}

		v.Set(slice)

	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)

			if field.PkgPath != "" {
				continue
			}

			var err error
			b, err = decode(b, v.Field(i))
			if err != nil {
				return nil, err
			}
		}

	default:
		return nil, errors.New("unsupported type in decodeList")
	}

	return rest, nil
}

func bytesToInt(b []byte) uint64 {
	if len(b) == 0 {
		return 0
	}
	return new(big.Int).SetBytes(b).Uint64()
}
