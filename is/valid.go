// Package is provides utility functions to validate various types and values.
package is

import (
	"reflect"

	"github.com/kubex-ecosystem/kbx/types"
	gl "github.com/kubex-ecosystem/logz"
)

func LogEntry(obj any) bool {
	if !Valid(obj) {
		return false
	}
	_, ok := obj.(gl.Entry)
	return ok
}

func Valid(obj any) bool {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return false
		}
		if v.Kind() == reflect.Ptr {
			if v.Elem().Kind() == reflect.Ptr && v.Elem().IsNil() {
				return false
			}
			v = v.Elem()
		}
	}
	if _, ok := types.KindMap[v.Kind()]; !ok {
		return false
	}
	if !v.IsValid() {
		return false
	}
	if v.IsZero() {
		return false
	}
	if v.Kind() == reflect.String && v.Len() == 0 {
		return false
	}
	if (v.Kind() == reflect.Slice || v.Kind() == reflect.Map || v.Kind() == reflect.Array) && v.Len() == 0 {
		return false
	}
	if v.Kind() == reflect.Bool {
		return true
	}
	return true
}

func Safe(obj any, strict bool) bool {
	v := reflect.ValueOf(obj)

	// nil pointers or invalid values
	if !v.IsValid() {
		return false
	}
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return false
		}
		v = v.Elem()
	}

	// zero value check (different meaning in strict vs resilient mode)
	if v.IsZero() {
		if strict {

			switch v.Kind() {
			case reflect.Bool, reflect.Int, reflect.Int64, reflect.Float64, reflect.String:
				// 0, false, "" são válidos em modo estrito
				return true
			}
		}
		return false
	}

	// empty collections → false no resilient mode
	switch v.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:
		if v.Len() == 0 {
			return !strict
		}
	}

	return true
}

func SpecVar(c uint8) bool {
	switch c {
	case '*', '#', '$', '@', '!', '?', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	}
	return false
}
func AlphaN(c uint8) bool {
	return c == '_' || '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}

func Alpha(c uint8) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_'
}

func Numeric(c any) bool {
	switch c := c.(type) {
	case uint8:
		return '0' <= c && c <= '9'
	default:
		return false
	}
}

func SameType[T any](obj any) bool {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return false
		}
		v = v.Elem()
	}
	if v.Type() != reflect.TypeFor[T]() {
		return false
	}
	return true
}

func Slice[T any](obj any) bool {
	v := reflect.ValueOf(obj)
	if SameType[any](obj) || !Valid(obj) {
		return false
	}
	if !types.KindMap[v.Kind()] && SameType[T](v.Interface()) {
		return false
	}
	if v.Type().Elem() != reflect.TypeFor[T]() {
		return false
	}
	return true
}

func Map[K, V any](obj any) bool {
	v := reflect.ValueOf(obj)
	if SameType[any](obj) || !Valid(obj) {
		return false
	}
	if !types.KindMap[v.Kind()] && (SameType[K](v.Interface()) || SameType[V](v.Interface())) {
		return false
	}
	if v.Type().Key() != reflect.TypeFor[K]() || v.Type().Elem() != reflect.TypeFor[V]() {
		return false
	}
	return true
}

func Struct[T any](obj any) bool {
	v := reflect.ValueOf(obj)
	if SameType[any](obj) || !Valid(obj) {
		return false
	}
	if !types.KindMap[v.Kind()] && SameType[T](v.Interface()) {
		return false
	}
	if v.Type() != reflect.TypeFor[T]() {
		return false
	}
	return true
}

func Compatible[T any](obj any) bool {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return false
		}
		v = v.Elem()
	}
	if !v.Type().ConvertibleTo(reflect.TypeFor[T]()) {
		return false
	}
	return true
}

func Implements[T any](obj any) bool {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return false
		}
		v = v.Elem()
	}
	if !v.Type().Implements(reflect.TypeFor[T]()) {
		return false
	}
	return true
}

func NilPtr(obj any) bool {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return true
		}
	}
	return false
}
