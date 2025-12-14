// Package get contains utility functions for retrieving values with defaults and type handling.
package get

import (
	"os"
	"reflect"

	gl "github.com/kubex-ecosystem/logz"
)

// EnvOr retrieves the value of the environment variable named by the key.
// If the variable is empty or not set, it returns the provided default value d.

func EnvOr(key, d string) string {
	value := os.Getenv(key)
	if value == "" {
		return d
	}
	return value
}

// Integer pointer helper functions - Convenience wrappers around generic Ptr function with type conversion

func UintPtr(n uint64) *uint { return Ptr(uint(n)) }
func IntPtr(n int64) *int    { return Ptr(int(n)) }

// Pointer helper functions - Convenience wrappers around generic Ptr function

func Uint64Ptr(n uint64) *uint64 { return Ptr(n) }
func Int64Ptr(n int64) *int64    { return Ptr(n) }
func BlPtr(b bool) *bool         { return Ptr(b) }
func Fl64Ptr(n float64) *float64 { return Ptr(n) }
func StrPtr(s string) *string    { return Ptr(s) }

// Generic Ptr function - Returns a pointer to the given value

func Ptr[T any](v T) *T { return &v }

// Type functions - Retrieve the reflect.Type and type name of a given value

func Type[T any](v T) reflect.Type { return reflect.TypeFor[T]() }
func TypeName(obj any) string {
	if t := Type(obj); t.Kind() == reflect.Pointer {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}


func SeedFromEnvMap[T any](prefix string, keyMap map[string]T, defMap map[string]T, ctlChan chan any) map[string]T {
	defer func(hCtl chan any) {
		if r := recover(); r != nil {
			// Handle the panic (e.g., log the error)
			gl.Errorf("Panic at the Hydration: %v", r)
			if ctlChan != nil {
				gl.Info("Async hydration due to panic recovery")
				for key, defaultValue := range defMap {
					keyMap[key] = ValOrType(keyMap[key], defaultValue)
				}
				ctlChan <- r
				return
			}
		}
	}(ctlChan)
	for key, defaultValue := range defMap {
		keyMap[key] = EnvOrType(prefix+"_"+key,
			ValOrType(keyMap[key], defaultValue),
		)
	}
	gl.Debugf("Hydrated Map for DBType %s: %+v", prefix, keyMap)
	return keyMap
}