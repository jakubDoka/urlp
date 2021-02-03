package urlp

import (
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/jakubDoka/sterr"
)

// Errors related to url Parsing
var (
	ErrNotPtr          = sterr.New("value is not pointer")
	ErrMissingTagValue = sterr.New("tag(%s) on %s is missing value")
	ErrMissingValue    = sterr.New("value with name %s is missing and is not optional")
	ErrInvalidField    = sterr.New("trying to parse %s into %s where type does not match")
	ErrNotSupported    = sterr.New("Type(%s) is not supported")
	ErrParseFail       = sterr.New("failed to parse %s to %s to field %s")
)

// Parse takes an url.Values and struct that is fed with them, you can use
// "optional" tag to make struct field optional or "name:"something"" to use custom name
func Parse(values url.Values, value interface{}) (err error) {
	ptr := reflect.ValueOf(value)
	if ptr.Kind() != reflect.Ptr {
		return ErrNotPtr
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

		f.init()

		err = f.set(values)
		if err != nil {
			return
		}
	}

	return nil
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
		return nil
	}

	val := vals[0]

	if f.Kind() == reflect.String {
		f.SetString(val)
		return
	}

	e := ErrParseFail.Args(val, f.Kind(), f.Name)
	switch f.Kind() {
	default:
		return ErrNotSupported.Args(f.Kind())
	case reflect.Bool:
		switch val {
		case "true":
			f.SetBool(true)
		case "false":
			f.SetBool(false)
		default:
			return e
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return e.Wrap(err)
		}
		f.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		n, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return e.Wrap(err)
		}
		f.SetUint(n)
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(val, f.StructField.Type.Bits())
		if err != nil {
			return e.Wrap(err)
		}
		f.SetFloat(n)
	}
	return nil
}
