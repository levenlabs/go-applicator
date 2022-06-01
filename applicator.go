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
	if val.Kind() == reflect.Invalid {
		return ErrCannotApply
	}
	var lastPtr reflect.Value
	var typ reflect.Type
	for {
		if val.Kind() == reflect.Invalid {
			return ErrCannotApply
		}
		typ = val.Type()

		if val.Kind() == reflect.Ptr {
			if val.IsNil() {
				return ErrCannotApply
			}
			lastPtr = val
			val = val.Elem()
		} else if val.Kind() == reflect.Interface {
			if val.IsNil() {
				return ErrCannotApply
			}
			val = val.Elem()
		} else {
			break
		}
	}
	switch typ.Kind() {
	case reflect.Struct:
		// if we didn't get a pointer at all then we can't apply
		if !lastPtr.IsValid() {
			return ErrCannotApply
		}
		val, err := h.applyStruct(val)
		if err != nil {
			return err
		}
		lastPtr.Elem().Set(val)
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
		// if we have an array of structs/arrays then we need a pointer to the array
		// otherwise we can't get the address
		// but if we have an array of interface{} then we don't know what to expect
		// but we should try to address them
		var shouldAddr bool
		var resetArray bool
		switch typ.Elem().Kind() {
		case reflect.Struct, reflect.Array:
			if !lastPtr.IsValid() {
				return ErrCannotApply
			}
			shouldAddr = true
			if !val.CanAddr() {
				resetArray = true
				ptr := reflect.New(val.Type())
				ptr.Elem().Set(val)
				val = ptr.Elem()
			}
		case reflect.Interface:
			shouldAddr = true
		}
		for i := 0; i < val.Len(); i++ {
			f := val.Index(i)
			if f.CanInterface() {
				if shouldAddr && f.CanAddr() {
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
		if resetArray {
			lastPtr.Elem().Set(val)
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

func (h *Applicator) applyStruct(strct reflect.Value) (reflect.Value, error) {
	if strct.Kind() != reflect.Struct {
		return strct, ErrCannotApply
	}
	if !strct.CanAddr() {
		ptr := reflect.New(strct.Type())
		ptr.Elem().Set(strct)
		strct = ptr.Elem()
	}
	typ := strct.Type()
	for i := 0; i < strct.NumField(); i++ {
		f := strct.Field(i)
		fn := strings.Split(typ.Field(i).Tag.Get(h.TagName), ",")
		if !f.CanSet() || fn[0] == "-" {
			continue
		}
		// if the field is a struct/slice/map/etc, run the applications on that
		// struct's field but continue if there was no errors in case there are
		// applications on the struct itself
		// we ignore ErrCannotApply since this is purely opportunistic
		if canApply(f.Type()) && f.CanInterface() {
			if f.CanAddr() && f.Kind() == reflect.Struct {
				if err := h.Apply(f.Addr().Interface()); err != nil && err != ErrCannotApply {
					return strct, err
				}
			} else {
				if err := h.Apply(f.Interface()); err != nil && err != ErrCannotApply {
					return strct, err
				}
			}
		}
		if fn[0] == "" {
			continue
		}
		fi := f.Interface()
		var err error
		for _, n := range fn {
			np := strings.SplitN(n, "=", 2)
			var no string
			if len(np) == 2 {
				no = np[1]
				n = np[0]
			}
			fn, ok := h.funcs[n]
			if !ok {
				return strct, ErrNotFound
			}

			fi, err = fn(fi, no)
			if err != nil {
				return strct, err
			}
		}
		if reflect.TypeOf(fi) != f.Type() {
			return strct, ErrInvalidSet
		}
		f.Set(reflect.ValueOf(fi))
	}
	return strct, nil
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
