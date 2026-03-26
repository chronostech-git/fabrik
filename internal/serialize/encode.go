package serialize

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"
)

func encodeValue(buf *bytes.Buffer, v reflect.Value) error {
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return binary.Write(buf, binary.BigEndian, uint8(0))
		}
		if err := binary.Write(buf, binary.BigEndian, uint8(1)); err != nil {
			return err
		}
		return encodeValue(buf, v.Elem())
	}

	switch v.Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Bool:
		return binary.Write(buf, binary.BigEndian, v.Interface())

	case reflect.String:
		b := []byte(v.String())
		if err := binary.Write(buf, binary.BigEndian, uint64(len(b))); err != nil {
			return err
		}
		buf.Write(b)
		return nil

	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			b := v.Bytes()
			if err := binary.Write(buf, binary.BigEndian, uint64(len(b))); err != nil {
				return err
			}
			buf.Write(b)
			return nil
		}

		length := v.Len()
		if err := binary.Write(buf, binary.BigEndian, uint64(length)); err != nil {
			return err
		}

		for i := 0; i < length; i++ {
			if err := encodeValue(buf, v.Index(i)); err != nil {
				return err
			}
		}
		return nil

	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if err := encodeValue(buf, v.Index(i)); err != nil {
				return err
			}
		}
		return nil

	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fieldType := v.Type().Field(i)
			fieldVal := v.Field(i)

			// 🔥 skip unexported fields
			if fieldType.PkgPath != "" {
				continue
			}

			// 🔥 skip interface fields (like Database)
			if fieldVal.Kind() == reflect.Interface {
				continue
			}

			if err := encodeValue(buf, fieldVal); err != nil {
				return err
			}

			if v.Kind() == reflect.Interface {
				if v.IsNil() {
					return binary.Write(buf, binary.BigEndian, uint8(0))
				}

				// mark non-nil
				if err := binary.Write(buf, binary.BigEndian, uint8(1)); err != nil {
					return err
				}

				return encodeValue(buf, v.Elem())
			}
		}
		return nil

	default:
		return errors.New("unsupported type: " + v.Kind().String())
	}
}
