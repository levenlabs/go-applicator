package helper

import (
	"reflect"
	"strings"
)

type function func(interface{}, string) (interface{}, error)

// Helper is an instance of
type Helper struct {
	funcs map[string]function
}

var defaultHelper = NewHelper()

// Returns a new instance of Helper with the builtin helpers added
func NewHelper() *Helper {
	return &Helper{
		funcs: map[string]function{
			"trim":  trim,
			"lower": lower,
		},
	}
}

// Runs all the helpers on the struct's fields. Must receive a pointer
func (h *Helper) Run(s interface{}) error {
	el := reflect.TypeOf(s)
	val := reflect.ValueOf(s)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return ErrUnsupported
	}
	el = el.Elem()
	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return ErrUnsupported
	}
	for i := 0; i < val.NumField(); i++ {
		fVal := val.Field(i)
		fEl := el.Field(i)
		fn := strings.Split(fEl.Tag.Get("helper"), ",")
		if !fVal.CanSet() || fn[0] == "" || fn[0] == "-" {
			continue
		}
		val := fVal.Interface()
		var err error
		for _, n := range fn {
			np := strings.SplitN(n, "=", 2)
			var no string
			if len(np) == 2 {
				no = np[1]
				n = np[0]
			}
			f, ok := h.funcs[n]
			if !ok {
				return ErrNotFound
			}

			val, err = f(val, no)
			if err != nil {
				return err
			}
		}
		if reflect.TypeOf(val) != fVal.Type() {
			return ErrInvalidSet
		}
		fVal.Set(reflect.ValueOf(val))
	}
	return nil
}

// Adds a new helper function
func (h *Helper) AddFunc(name string, f function) {
	h.funcs[name] = f
}

// Runs the default helper
func Run(s interface{}) error {
	return defaultHelper.Run(s)
}

// Adds a new helper function to the default helper
func AddFunc(name string, f function) {
	defaultHelper.AddFunc(name, f)
}
