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
			"trim":    trim,
			"lower":   lower,
			"fillNil": fillNil,
		},
	}
}

func canApply(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Struct, reflect.Slice, reflect.Array, reflect.Map, reflect.Interface:
		return true
	case reflect.Ptr:
		return canApply(typ.Elem())
	default:
		return false
	}
}

// Apply runs all the functions on the struct's fields based on their tags.
// Accepts a pointer to a struct, a map, slice, pointer to an array, or interface
// to any of the previous types.
func (h *Applicator) Apply(s interface{}) error {
	val := reflect.ValueOf(s)
	typ := val.Type()
	lastPtr := val
	for {
		if val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
			if val.IsNil() {
				return ErrCannotApply
			}
			lastPtr = val
			val = val.Elem()
			typ = val.Type()
		} else {
			break
		}
	}
	switch typ.Kind() {
	case reflect.Struct:
		return h.applyStruct(lastPtr)
	case reflect.Slice:
		// make sure its a slice of compatible types
		if !canApply(typ.Elem()) {
			return ErrCannotApply
		}
		for i := 0; i < val.Len(); i++ {
			f := val.Index(i)
			if f.CanInterface() {
				if f.Kind() == reflect.Struct {
					f = f.Addr()
				}
				if err := h.Apply(f.Interface()); err != nil {
					// we ignore this error since it might be a []interface{} and some of them
					// might not be compatible but some might be
					if err == ErrCannotApply {
						continue
					}
					return err
				}
			}
		}
	case reflect.Array:
		// make sure its an array of compatible types
		if !canApply(typ.Elem()) {
			return ErrCannotApply
		}
		for i := 0; i < val.Len(); i++ {
			f := val.Index(i)
			if f.CanInterface() {
				if f.Kind() == reflect.Struct {
					// this means we didn't get an addressable Array passed
					if !f.CanAddr() {
						return ErrCannotApply
					}
					f = f.Addr()
				}
				if err := h.Apply(f.Interface()); err != nil {
					// we ignore this error since it might be a []interface{} and some of them
					// might not be compatible but some might be
					if err == ErrCannotApply {
						continue
					}
					return err
				}
			}
		}
	case reflect.Map:
		// make sure its a map of compatible types
		if !canApply(typ.Elem()) {
			return ErrCannotApply
		}
		for _, k := range val.MapKeys() {
			f := val.MapIndex(k)
			if f.CanInterface() {
				if f.Kind() == reflect.Struct {
					// this means the values aren't addressable
					if !f.CanAddr() {
						return ErrCannotApply
					}
					f = f.Addr()
				}
				if err := h.Apply(f.Interface()); err != nil {
					// we ignore this error since it might be a []interface{} and some of them
					// might not be compatible but some might be
					if err == ErrCannotApply {
						continue
					}
					return err
				}
			}
		}
	default:
		return ErrCannotApply
	}
	return nil
}

func (h *Applicator) applyStruct(val reflect.Value) error {
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return ErrCannotApply
	}
	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return ErrCannotApply
	}
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		fVal := val.Field(i)
		fn := strings.Split(typ.Field(i).Tag.Get(h.TagName), ",")
		if !fVal.CanSet() || fn[0] == "-" {
			continue
		}
		// if the field is a struct/slice/map/etc, run the applications on that
		// struct's field but continue if there was no errors in case there are
		// applications on the struct itself
		// we ignore ErrCannotApply since this is purely opportunistic
		if canApply(fVal.Type()) && fVal.CanInterface() {
			if fVal.CanAddr() && fVal.Kind() == reflect.Struct {
				if err := h.Apply(fVal.Addr().Interface()); err != nil && err != ErrCannotApply {
					return err
				}
			} else {
				if err := h.Apply(fVal.Interface()); err != nil && err != ErrCannotApply {
					return err
				}
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
