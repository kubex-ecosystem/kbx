// Package is provides utility functions to validate various types and values.
package is

import (
	"reflect"

	"github.com/kubex-ecosystem/kbx/types"
	gl "github.com/kubex-ecosystem/logz"
)

// LogEntry checks if the given object is of type gl.Entry
func LogEntry(obj any) bool {
	if !Valid(obj) {
		return false
	}
	_, ok := obj.(gl.Entry)
	return ok
}

// Valid checks if the given object is valid (not nil, not zero value, etc.)
// This function checks like JS, Python truthy/falsy values.
// It returns false for nil pointers, zero values, empty strings, and empty collections.
// Its behavior is always resilient and very strict.
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

// PtrOf checks if the given object is a non-nil pointer to type T.
func PtrOf[T any](obj any) bool {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Pointer {
		return false
	}
	if v.IsNil() {
		return false
	}
	if v.Elem().Type() != reflect.TypeFor[T]() {
		return false
	}
	return true
}

// Safe checks if the given object is "safe" to use. Its seems like Valid, but with different rules.
// In strict mode, it considers zero values of basic types (0, false, "") as safe.
// In resilient mode, it treats empty collections as unsafe.
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

// SpecVar checks if the given character is a special variable character.
// Special variable characters include: '*', '#', '$', '@', '!', '?', '-', and digits '0'-'9'.
func SpecVar(c uint8) bool {
	switch c {
	case '*', '#', '$', '@', '!', '?', '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	}
	return false
}

// AlphaN checks if the given character is alphanumeric or an underscore.
// It returns true for characters 'a'-'z', 'A'-'Z', '0'-'9', and '_'.
func AlphaN(c uint8) bool {
	return c == '_' || '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}

// Alpha checks if the given character is an alphabetic character or an underscore.
// It returns true for characters 'a'-'z', 'A'-'Z', and '_'.
func Alpha(c uint8) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_'
}

// Numeric checks if the given character is a numeric digit.
// It returns true for characters '0'-'9'.
func Numeric(c any) bool {
	switch c := c.(type) {
	case uint8:
		return '0' <= c && c <= '9'
	default:
		return false
	}
}

// SameType checks if the given object is of the same type as T.
// It returns false for nil pointers or mismatched types.
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

// Slice checks if the given object is a slice of type T.
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

// Map checks if the given object is a map with key type K and value type V.
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

// Struct checks if the given object is a struct of type T.
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

// Compatible checks if the given object is convertible to type T.
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

// Compatible checks if the given object is convertible to type T.
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

// NilPtr checks if the given object is a nil pointer or interface
func NilPtr(obj any) bool {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return true
		}
	}
	return false
}

// ArrayObj checks if the object o exists in the array a of type T.
func ArrayObj[T any](o T, a []T) bool {
	for _, v := range a {
		if reflect.DeepEqual(o, v) {
			return true
		}
	}
	return false
}
