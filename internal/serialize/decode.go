package serialize

import (
	"encoding/binary"
	"errors"
	"io"
	"reflect"
)

func decodeValue(r io.Reader, v reflect.Value) error {
	if v.Kind() == reflect.Pointer {
		var flag uint8
		if err := binary.Read(r, binary.BigEndian, &flag); err != nil {
			return err
		}
		if flag == 0 {
			v.Set(reflect.Zero(v.Type()))
			return nil
		}
		v.Set(reflect.New(v.Type().Elem()))
		return decodeValue(r, v.Elem())
	}

	switch v.Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Bool:
		return binary.Read(r, binary.BigEndian, v.Addr().Interface())

	case reflect.String:
		var length uint64
		if err := binary.Read(r, binary.BigEndian, &length); err != nil {
			return err
		}
		b := make([]byte, length)
		if _, err := io.ReadFull(r, b); err != nil {
			return err
		}
		v.SetString(string(b))
		return nil

	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			var length uint64
			if err := binary.Read(r, binary.BigEndian, &length); err != nil {
				return err
			}
			b := make([]byte, length)
			if _, err := io.ReadFull(r, b); err != nil {
				return err
			}
			v.SetBytes(b)
			return nil
		}

		var length uint64
		if err := binary.Read(r, binary.BigEndian, &length); err != nil {
			return err
		}

		slice := reflect.MakeSlice(v.Type(), int(length), int(length))
		for i := 0; i < int(length); i++ {
			if err := decodeValue(r, slice.Index(i)); err != nil {
				return err
			}
		}
		v.Set(slice)
		return nil

	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if err := decodeValue(r, v.Index(i)); err != nil {
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

			// 🔥 skip interface fields
			if fieldVal.Kind() == reflect.Interface {
				continue
			}

			if err := decodeValue(r, fieldVal); err != nil {
				return err
			}

			if v.Kind() == reflect.Interface {
				var flag uint8
				if err := binary.Read(r, binary.BigEndian, &flag); err != nil {
					return err
				}

				if flag == 0 {
					v.Set(reflect.Zero(v.Type()))
					return nil
				}

				// ⚠️ we cannot instantiate interface directly
				// so just skip for now (or error)

				return errors.New("cannot decode into interface without concrete type")
			}
		}
		return nil

	default:
		return errors.New("unsupported type: " + v.Kind().String())
	}
}
