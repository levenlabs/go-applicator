package helper

import (
	"reflect"
	"strings"
)

func trim(i interface{}, _ string) (interface{}, error) {
	t := func(v reflect.Value) string {
		return strings.TrimSpace(v.String())
	}

	v := reflect.ValueOf(i)
	var err error
	switch v.Kind() {
	case reflect.String:
		i = t(v)
	case reflect.Ptr:
		v = v.Elem()
		if v.Kind() != reflect.String {
			err = ErrUnsupported
		} else {
			str := t(v)
			i = &str
		}
	default:
		err = ErrUnsupported
	}
	return i, err
}

func lower(i interface{}, _ string) (interface{}, error) {
	l := func(v reflect.Value) string {
		return strings.ToLower(v.String())
	}

	v := reflect.ValueOf(i)
	var err error
	switch v.Kind() {
	case reflect.String:
		i = l(v)
	case reflect.Ptr:
		v = v.Elem()
		if v.Kind() != reflect.String {
			err = ErrUnsupported
		} else {
			str := l(v)
			i = &str
		}
	default:
		err = ErrUnsupported
	}
	return i, err
}
