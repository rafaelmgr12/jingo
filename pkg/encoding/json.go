package encoding

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/rafaelmgr12/jingo/pkg/parser"
)

// Marshal converts a Go value into a JSON string with optional configuration.
// It handles all basic Go types including interface{}, maps, slices, arrays, and structs.
func Marshal(v interface{}, opts ...Option) ([]byte, error) {
	options, err := applyOptions(opts...)
	if err != nil {
		return nil, err
	}

	value, err := marshalValue(reflect.ValueOf(v))
	if err != nil {
		return nil, fmt.Errorf("marshal error: %v", err)
	}

	var b strings.Builder
	if err := writeValue(&b, value); err != nil {
		return nil, fmt.Errorf("writing error: %v", err)
	}

	result := []byte(b.String())
	if len(result) > options.MaxSize {
		return nil, fmt.Errorf("marshaled JSON size (%d bytes) exceeds maximum allowed size (%d bytes)",
			len(result), options.MaxSize)
	}

	return result, nil
}

// Unmarshal parses JSON data and stores the result in the value pointed to by v.
// The target value must be a non-nil pointer.
func Unmarshal(data []byte, v interface{}, opts ...Option) error {
	options, err := applyOptions(opts...)
	if err != nil {
		return err
	}

	if len(data) > options.MaxSize {
		return fmt.Errorf("input JSON size (%d bytes) exceeds maximum allowed size (%d bytes)",
			len(data), options.MaxSize)
	}

	l := parser.NewLexer(string(data))
	p := parser.NewParser(l)

	value, err := p.ParseJSON()
	if err != nil {
		return fmt.Errorf("parse error: %v", err)
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("unmarshal target must be a non-nil pointer")
	}

	return unmarshalValue(value, rv.Elem())
}

// marshalValue converts a reflect.Value to a parser.Value
func marshalValue(v reflect.Value) (parser.Value, error) {
	// Handle interface{} values
	if v.Kind() == reflect.Interface && !v.IsNil() {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.String:
		return &parser.StringLiteral{
			Value: v.String(),
			Token: parser.Token{Type: parser.TokenString},
		}, nil

	case reflect.Bool:
		return &parser.Boolean{
			Value: v.Bool(),
			Token: parser.Token{Type: parser.TokenTrue},
		}, nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		num := parser.NewNumberLiteral(parser.Token{
			Type:    parser.TokenNumber,
			Literal: fmt.Sprintf("%d", v.Int()),
		})

		return num, nil

	case reflect.Float32, reflect.Float64:
		num := parser.NewNumberLiteral(parser.Token{
			Type:    parser.TokenNumber,
			Literal: fmt.Sprintf("%g", v.Float()),
		})

		return num, nil

	case reflect.Map:
		if v.Type().Key().Kind() != reflect.String {
			return nil, fmt.Errorf("map key must be string")
		}

		obj := &parser.Object{
			Token: parser.Token{Type: parser.TokenBraceOpen},
			Pairs: make(map[string]parser.Value),
		}

		iter := v.MapRange()
		for iter.Next() {
			value, err := marshalValue(iter.Value())
			if err != nil {
				return nil, fmt.Errorf("map value: %v", err)
			}

			obj.Pairs[iter.Key().String()] = value
		}

		return obj, nil

	case reflect.Slice, reflect.Array:
		arr := &parser.Array{
			Token:    parser.Token{Type: parser.TokenBracketOpen},
			Elements: make([]parser.Value, 0, v.Len()),
		}

		for i := 0; i < v.Len(); i++ {
			value, err := marshalValue(v.Index(i))
			if err != nil {
				return nil, fmt.Errorf("index %d: %v", i, err)
			}

			arr.Elements = append(arr.Elements, value)
		}

		return arr, nil

	case reflect.Ptr:
		if v.IsNil() {
			return &parser.Null{Token: parser.Token{Type: parser.TokenNull}}, nil
		}

		return marshalValue(v.Elem())

	case reflect.Struct:
		obj := &parser.Object{
			Token: parser.Token{Type: parser.TokenBraceOpen},
			Pairs: make(map[string]parser.Value),
		}

		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)

			tag := field.Tag.Get("json")
			if tag == "-" {
				continue
			}

			name := field.Name

			if tag != "" {
				tagParts := strings.Split(tag, ",")
				if len(tagParts) > 0 && tagParts[0] != "" {
					name = tagParts[0]
				}
			}

			value, err := marshalValue(v.Field(i))
			if err != nil {
				return nil, fmt.Errorf("field %s: %v", name, err)
			}

			obj.Pairs[name] = value
		}

		return obj, nil

	case reflect.Interface:
		if v.IsNil() {
			return &parser.Null{Token: parser.Token{Type: parser.TokenNull}}, nil
		}

		return marshalValue(v.Elem())

	default:
		return nil, fmt.Errorf("unsupported type: %v", v.Type())
	}
}

// unmarshalValue converts a parser.Value to a reflect.Value
func unmarshalValue(v parser.Value, rv reflect.Value) error {
	if v == nil {
		return fmt.Errorf("cannot unmarshal nil value")
	}

	if rv.Kind() == reflect.Interface && rv.NumMethod() == 0 {
		switch val := v.(type) {
		case *parser.Object:
			obj := map[string]interface{}{}

			for k, v := range val.Pairs {
				var mapValue interface{}
				if err := unmarshalValue(v, reflect.ValueOf(&mapValue).Elem()); err != nil {
					return fmt.Errorf("map key %q: %v", k, err)
				}

				obj[k] = mapValue
			}

			rv.Set(reflect.ValueOf(obj))

		case *parser.Array:
			arr := make([]interface{}, len(val.Elements))

			for i, elem := range val.Elements {
				var arrayValue interface{}
				if err := unmarshalValue(elem, reflect.ValueOf(&arrayValue).Elem()); err != nil {
					return fmt.Errorf("index %d: %v", i, err)
				}

				arr[i] = arrayValue
			}

			rv.Set(reflect.ValueOf(arr))

		case *parser.StringLiteral:
			rv.Set(reflect.ValueOf(val.Value))

		case *parser.NumberLiteral:
			if val.IsInt {
				rv.Set(reflect.ValueOf(val.Int))
			} else {
				rv.Set(reflect.ValueOf(val.Float))
			}

		case *parser.Boolean:
			rv.Set(reflect.ValueOf(val.Value))

		case *parser.Null:
			rv.Set(reflect.Zero(rv.Type()))

		default:
			return fmt.Errorf("unknown value type: %T", v)
		}

		return nil
	}

	switch val := v.(type) {
	case *parser.Object:
		return unmarshalObject(val, rv)

	case *parser.Array:
		return unmarshalArray(val, rv)

	case *parser.StringLiteral:
		return unmarshalString(val, rv)

	case *parser.NumberLiteral:
		return unmarshalNumber(val, rv)

	case *parser.Boolean:
		return unmarshalBool(val, rv)

	case *parser.Null:
		return unmarshalNull(rv)

	default:
		return fmt.Errorf("unknown value type: %T", v)
	}
}

// unmarshalObject handles unmarshaling of JSON objects into Go structs or maps
func unmarshalObject(obj *parser.Object, rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.Map:
		if rv.IsNil() {
			rv.Set(reflect.MakeMap(rv.Type()))
		}

		for k, v := range obj.Pairs {
			elemType := rv.Type().Elem()
			mapValue := reflect.New(elemType).Elem()

			if err := unmarshalValue(v, mapValue); err != nil {
				return fmt.Errorf("map value %q: %v", k, err)
			}

			rv.SetMapIndex(reflect.ValueOf(k), mapValue)
		}

	case reflect.Struct:
		t := rv.Type()

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)

			tag := field.Tag.Get("json")
			if tag == "-" {
				continue
			}

			name := field.Name
			if tag != "" {
				name = strings.Split(tag, ",")[0]
			}

			if v, ok := obj.Pairs[name]; ok {
				if err := unmarshalValue(v, rv.Field(i)); err != nil {
					return fmt.Errorf("field %s: %v", name, err)
				}
			}
		}

	default:
		return fmt.Errorf("cannot unmarshal object into %v", rv.Type())
	}

	return nil
}

// unmarshalArray handles unmarshaling of JSON arrays into Go slices or arrays
func unmarshalArray(arr *parser.Array, rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.Slice:
		slice := reflect.MakeSlice(rv.Type(), len(arr.Elements), len(arr.Elements))
		for i, elem := range arr.Elements {
			if err := unmarshalValue(elem, slice.Index(i)); err != nil {
				return fmt.Errorf("index %d: %v", i, err)
			}
		}

		rv.Set(slice)

	case reflect.Array:
		if rv.Len() != len(arr.Elements) {
			return fmt.Errorf("cannot unmarshal array of length %d into array of length %d",
				len(arr.Elements), rv.Len())
		}

		for i, elem := range arr.Elements {
			if err := unmarshalValue(elem, rv.Index(i)); err != nil {
				return fmt.Errorf("index %d: %v", i, err)
			}
		}

	default:
		return fmt.Errorf("cannot unmarshal array into %v", rv.Type())
	}

	return nil
}

// unmarshalString handles unmarshaling of JSON strings into Go strings
func unmarshalString(str *parser.StringLiteral, rv reflect.Value) error {
	if rv.Kind() != reflect.String {
		return fmt.Errorf("cannot unmarshal string into %v", rv.Type())
	}

	rv.SetString(str.Value)

	return nil
}

// unmarshalNumber handles unmarshaling of JSON numbers into Go numeric types
func unmarshalNumber(num *parser.NumberLiteral, rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if !num.IsInt {
			return fmt.Errorf("cannot unmarshal float into %v", rv.Type())
		}

		rv.SetInt(num.Int)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if !num.IsInt || num.Int < 0 {
			return fmt.Errorf("cannot unmarshal negative number into %v", rv.Type())
		}

		rv.SetUint(uint64(num.Int))

	case reflect.Float32, reflect.Float64:
		rv.SetFloat(num.Float)

	default:
		return fmt.Errorf("cannot unmarshal number into %v", rv.Type())
	}

	return nil
}

// unmarshalBool handles unmarshaling of JSON booleans into Go bools
func unmarshalBool(b *parser.Boolean, rv reflect.Value) error {
	if rv.Kind() != reflect.Bool {
		return fmt.Errorf("cannot unmarshal boolean into %v", rv.Type())
	}

	rv.SetBool(b.Value)

	return nil
}

// unmarshalNull handles unmarshaling of JSON null into Go values
func unmarshalNull(rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.Interface, reflect.Ptr, reflect.Map, reflect.Slice:
		rv.Set(reflect.Zero(rv.Type()))
		return nil
	default:
		return fmt.Errorf("cannot unmarshal null into %v", rv.Type())
	}
}

// writeValue writes a parser.Value to a strings.Builder
func writeValue(b *strings.Builder, v parser.Value) error {
	switch val := v.(type) {
	case *parser.Object:
		b.WriteString("{")

		i := 0
		for k, v := range val.Pairs {
			if i > 0 {
				b.WriteString(",")
			}

			fmt.Fprintf(b, "%q:", k)

			if err := writeValue(b, v); err != nil {
				return err
			}

			i++
		}

		b.WriteString("}")

	case *parser.Array:
		b.WriteString("[")

		for i, v := range val.Elements {
			if i > 0 {
				b.WriteString(",")
			}

			if err := writeValue(b, v); err != nil {
				return err
			}
		}

		b.WriteString("]")

	case *parser.StringLiteral:
		fmt.Fprintf(b, "%q", val.Value)

	case *parser.NumberLiteral:
		b.WriteString(val.String())

	case *parser.Boolean:
		if val.Value {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}

	case *parser.Null:
		b.WriteString("null")

	default:
		return fmt.Errorf("unknown value type: %T", v)
	}

	return nil
}
