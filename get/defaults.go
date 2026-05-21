// Package get contains utility functions for retrieving values with defaults and type handling.
package get

import (
	"reflect"

	gl "github.com/kubex-ecosystem/logz"
)

// Integer pointer helper functions - Convenience wrappers around generic Ptr function with type conversion

// UintPtr returns a pointer to the given value n of type uint.
func UintPtr(n uint64) *uint { return Ptr(uint(n)) }

// IntPtr returns a pointer to the given value n of type int.
func IntPtr(n int64) *int { return Ptr(int(n)) }

// Pointer helper functions - Convenience wrappers around generic Ptr function

// Uint64Ptr returns a pointer to the given value n of type uint64.
func Uint64Ptr(n uint64) *uint64 { return Ptr(n) }

// Int64Ptr returns a pointer to the given value n of type int64.
func Int64Ptr(n int64) *int64 { return Ptr(n) }

// BlPtr returns a pointer to the given value b of type bool.
func BlPtr(b bool) *bool { return Ptr(b) }

// Fl64Ptr returns a pointer to the given value n of type float64.
func Fl64Ptr(n float64) *float64 { return Ptr(n) }

// StrPtr returns a pointer to the given value s of type string.
func StrPtr(s string) *string { return Ptr(s) }

// Generic Ptr function - Returns a pointer to the given value

// Ptr returns a pointer to the given value v of any type T.
func Ptr[T any](v T) *T { return new(v) }

// Type functions - Retrieve the reflect.Type and type name of a given value
func Type[T any](v T) reflect.Type { return reflect.TypeFor[T]() }

// TypeName returns the type name of a given value.
func TypeName(obj any) string {
	if t := Type(obj); t.Kind() == reflect.Pointer {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

// SeedFromEnvMap hydrates a map of string keys to values of type T
// by checking for environment variables with a specified prefix.
// For each key in defMap, it looks for an environment variable named
// prefix_key and sets the corresponding value in keyMap.
// If the environment variable is not set, it uses the default value from defMap.
// If ctlChan is provided, it will be used to signal completion or errors asynchronously.
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
