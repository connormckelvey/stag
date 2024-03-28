package tagram

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func parsePrimitive(literal string, target any) error {
	switch v := target.(type) {
	case *string:
		*v = literal
	case *complex64:
		c128, err := strconv.ParseComplex(literal, 64)
		if err != nil {
			return err
		}
		*v = complex64(c128)
	case *complex128:
		c128, err := strconv.ParseComplex(literal, 64)
		if err != nil {
			return err
		}
		*v = c128
	case *float32:
		f64, err := strconv.ParseFloat(literal, 32)
		if err != nil {
			return err
		}
		*v = float32(f64)
	case *float64:
		f64, err := strconv.ParseFloat(literal, 32)
		if err != nil {
			return err
		}
		*v = f64
	case *uint8:
		u64, err := strconv.ParseUint(literal, 10, 8)
		if err != nil {
			return err
		}
		*v = uint8(u64)
	case *uint16:
		u64, err := strconv.ParseUint(literal, 10, 16)
		if err != nil {
			return err
		}
		*v = uint16(u64)
	case *uint32:
		u64, err := strconv.ParseUint(literal, 10, 32)
		if err != nil {
			return err
		}
		*v = uint32(u64)
	case *uint64:
		u64, err := strconv.ParseUint(literal, 10, 8)
		if err != nil {
			return err
		}
		*v = u64
	case *int:
		i, err := strconv.Atoi(literal)
		if err != nil {
			return err
		}
		*v = i
	case *int8:
		i64, err := strconv.ParseInt(literal, 10, 8)
		if err != nil {
			return err
		}
		*v = int8(i64)
	case *int16:
		i64, err := strconv.ParseInt(literal, 10, 16)
		if err != nil {
			return err
		}
		*v = int16(i64)
	case *int32: // alias to rune
		i64, err := strconv.ParseInt(literal, 10, 32)
		if err != nil {
			return err
		}
		*v = int32(i64)
	case *int64:
		i64, err := strconv.ParseInt(literal, 10, 8)
		if err != nil {
			return err
		}
		*v = i64
	case *bool:
		b, err := strconv.ParseBool(literal)
		if err != nil {
			return err
		}
		*v = b
	default:
		return errors.New("unsupported primitive")
	}
	return nil
}

func parseInto(token string, targetFieldValue reflect.Value) error {
	dst := targetFieldValue.Addr().Interface()

	switch targetFieldValue.Kind() {
	case reflect.Slice:
		// a;b;c
		items := strings.Split(token, ";")
		var values []reflect.Value
		se := targetFieldValue.Type().Elem()
		for _, item := range items {
			v := reflect.New(se)
			err := parsePrimitive(item, v.Interface())
			if err != nil {
				return err
			}
			values = append(values, v.Elem())
		}
		sa := reflect.Append(targetFieldValue, values...)
		targetFieldValue.Set(sa)
	case reflect.Map:
		return fmt.Errorf("type 'map' not implemented")
	default:
		err := parsePrimitive(token, dst)
		if err != nil {
			return err
		}
	}
	return nil
}

func requireStruct(aStruct any) (reflect.Type, error) {
	structType := reflect.TypeOf(aStruct)
	if structType.Kind() == reflect.Pointer {
		structType = structType.Elem()
	}
	if structType.Kind() != reflect.Struct {
		return nil, errors.New("provided value not a struct type")
	}
	return structType, nil
}

func identity[T any]() (id T) {
	return
}
