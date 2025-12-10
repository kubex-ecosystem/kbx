// Package get contains utility functions for retrieving values with defaults and type handling.
package get

import (
	"os"
	"reflect"
)

func TypeName(obj any) string {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

func EnvOr(key, d string) string {
	value := os.Getenv(key)
	if value == "" {
		return d
	}
	return value
}

func BlPtr(b bool) *bool {
	return &b
}

func Fl64Ptr(n float64) *float64 {
	return &n
}

func IntPtr(n int64) *int {
	ni := int(n)
	return &ni
}

func StrPtr(s string) *string {
	return &s
}

func Type[T any](v T) reflect.Type {
	return reflect.TypeFor[T]()
}
