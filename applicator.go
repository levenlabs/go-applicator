package applicator

import (
	"reflect"
	"strings"
)

type function func(interface{}, string) (interface{}, error)

// Applicator is an instance of defined functions to run on fields in a struct
type Applicator struct {
	TagName string
	funcs   map[string]function
}

var defaultApplicator = New()

// New returns an instance of Applicator with the builtin functions added
// The default TagName is apply
func New() *Applicator {
	return &Applicator{
		TagName: "apply",
		funcs: map[string]function{
			"trim":   trim,
			"lower":  lower,
			"fillNil": fillNil,
		},
	}
}

// Apply runs all the functions on the struct's fields based on their tags.
// Must receive a struct pointer
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
		if !fVal.CanSet() || fn[0] == "-" {
			continue
		}
		// if the field is a struct, run the applications on that struct's field
		// but continue if there was no errors in case there are applications on
		// the struct itself
		if fVal.Kind() == reflect.Ptr && !fVal.IsNil() && fVal.Elem().Kind() == reflect.Struct {
			if err := h.Apply(fVal.Interface()); err != nil {
				return err
			}
		} else if fVal.Kind() == reflect.Struct {
			// we need to get a pointer to the struct
			if err := h.Apply(fVal.Addr().Interface()); err != nil {
				return err
			}
		}
		if fn[0] == "" {
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

// AddFunc adds a new function to this Applicator for the given name
// the function must have the definition:
// `func(interface{}, string) (interface{}, error)`
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
