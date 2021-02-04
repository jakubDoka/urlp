package urlp

import (
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/jakubDoka/sterr"
)

// errors
var (
	ErrNotPtr       = sterr.New("value is not a pointer but %s")
	ErrMissingValue = sterr.New("value with name %s is missing and does not have 'optional' tag")
	ErrNotSupported = sterr.New("type %s is not supported")
	ErrParseFail    = sterr.New("failed to parse %s into %s when filling field %s")
)

// Parse takes an url.Values and struct that is fed with them, you can use
// "optional" tag to make struct field optional or "something" to use custom name,
// you can nest struct but all fields will be considered at same level so:
//
//	type foo {
//		A, B, C int
//	}
//
// acts same as:
//
//	type bar {
//		A int
//		B struct {
//			B, C int
//		}
//	}
//
// `urlp:"name,optional"` will assing value under "name" key to tagged field and will not raise a
// error if field is missing
func Parse(values url.Values, value interface{}) (err error) {
	ptr := reflect.ValueOf(value)
	if ptr.Kind() != reflect.Ptr {
		return ErrNotPtr.Args(ptr.Kind())
	}

	v := ptr.Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := field{
			StructField: t.Field(i),
			Value:       v.Field(i),
		}

		if !f.CanSet() {
			continue
		}

		if f.CanAddr() && f.Kind() == reflect.Struct {
			err = Parse(values, f.Value.Addr().Interface())
			if err != nil {
				return
			}
			continue
		}

		f.init()

		err = f.set(values)
		if err != nil {
			return
		}
	}

	return
}

type field struct {
	reflect.StructField
	reflect.Value
	optional bool
}

func (f *field) init() {
	raw, ok := f.Tag.Lookup("urlp")
	if !ok {
		return
	}

	tags := strings.Split(raw, ",")
	for _, t := range tags {
		if t == "optional" {
			f.optional = true
		} else {
			f.Name = t
		}
	}
}

func (f *field) set(values url.Values) (err error) {
	vals, ok := values[f.Name]
	if !ok || len(vals) == 0 {
		if !f.optional {
			return ErrMissingValue.Args(f.Name)
		}
		return
	}

	if f.Kind() == reflect.Slice {
		slice := reflect.MakeSlice(f.StructField.Type, len(vals), len(vals))
		for i, val := range vals {
			err = setAny(val, f.Name+strconv.Itoa(i), false, f.StructField.Type.Elem(), slice.Index(i))
			if err != nil {
				return
			}
		}
		f.Set(slice)
		return
	}

	return setAny(vals[0], f.Name, f.optional, f.StructField.Type, f.Value)
}

func setAny(val, name string, optional bool, t reflect.Type, f reflect.Value) (err error) {
	var (
		k = f.Kind()
		e = ErrParseFail.Args(val, k, name)
	)

	switch k {
	case reflect.String:
		f.SetString(val)
		return
	case reflect.Bool:
		switch val {
		case "true":
			f.SetBool(true)
		case "false":
			f.SetBool(false)
		default:
			return e
		}
		return
	}

	if val == "" {
		if optional {
			return e
		}
		return
	}

	c := val[0]
	if c != '-' && (c < '0' || c > '9') {
		return e
	}

	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var n int64
		n, err = strconv.ParseInt(val, 10, 64)
		if err == nil {
			f.SetInt(n)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		var n uint64
		n, err = strconv.ParseUint(val, 10, 64)
		if err == nil {
			f.SetUint(n)
		}
	case reflect.Float32, reflect.Float64:
		var n float64
		n, err = strconv.ParseFloat(val, t.Bits())
		if err == nil {
			f.SetFloat(n)
		}
	default:
		return ErrNotSupported.Args(f.Kind())
	}

	return e.Wrap(err)
}
