package applicator

import (
	"reflect"
	"strings"
)

type function func(interface{}, string) (interface{}, error)

// Applicator is an instance of
type Applicator struct {
	TagName string
	funcs map[string]function
}

var defaultApplicator = New()

// New returns an instance of Applicator with the builtin functions added
func New() *Applicator {
	return &Applicator{
		TagName: "apply",
		funcs: map[string]function{
			"trim":  trim,
			"lower": lower,
		},
	}
}

// Runs all the functions on the struct's fields. Must receive a pointer
func (h *Applicator) Apply(s interface{}) error {
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
		fn := strings.Split(fEl.Tag.Get(h.TagName), ",")
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

// Adds a new function
func (h *Applicator) AddFunc(name string, f function) {
	h.funcs[name] = f
}

// Apply calls Apply on the default Applicator
func Apply(s interface{}) error {
	return defaultApplicator.Apply(s)
}

// AddFunc calls AddFunc on the default Applicator
func AddFunc(name string, f function) {
	defaultApplicator.AddFunc(name, f)
}
