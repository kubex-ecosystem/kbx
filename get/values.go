// Package get provides utility functions for retrieving values with default fallbacks and type handling.
// It includes functions to get values from environment variables, handle errors, and extract file extensions.
package get

import (
	"encoding/json"
	"os"
	"reflect"

	"github.com/kubex-ecosystem/kbx/is"
	gl "github.com/kubex-ecosystem/logz"
)

// ValueOr returns the value if it's valid,
// otherwise returns the default value along with its type.
func ValueOr[T any](value T, d T) (T, reflect.Type) {
	if !is.Valid(value) {
		return d, reflect.TypeFor[T]()
	}
	return value, reflect.TypeFor[T]()
}

// ValErrOr executes the provided function and returns
// its result if valid, otherwise returns the default value.
func ValErrOr[T any](fn func() (T, error), d T) T {
	if is.Safe(fn, false) {
		value, err := fn()
		if err != nil || !is.Valid(value) {
			gl.Errorf("ValErrOr[%s] failed: %v", reflect.TypeFor[T]().String(), err)
			return d
		}
		return value
	}
	gl.Errorf("We could not safely get the value (%s)", reflect.TypeFor[T]().String())
	return d
}

// ValOrType returns the value if it's valid,
// otherwise returns the default value.
func ValOrType[T any](value T, d T) T {
	if !is.Valid(value) {
		return d
	}
	return value
}

// EnvOrType retrieves an environment variable by key and attempts to convert it to the specified type T.
// If the environment variable is not set or conversion fails, it returns the default value.
func EnvOrType[T any](key string, d T) T {
	value := os.Getenv(key)
	// Sempre vem texto da env
	if len(value) == 0 {
		return d
	}
	if reflect.ValueOf(value).CanConvert(reflect.TypeFor[T]()) {
		return reflect.ValueOf(value).Convert(reflect.TypeFor[T]()).Interface().(T)
	}
	var result T
	if err := json.Unmarshal([]byte(value), &result); err != nil {
		return d
	}
	if is.Safe(result, false) {
		return result
	}
	return result
}

// ValOrAny returns the value if it's valid, otherwise returns the default value.
func ValOrAny[T any](value T, d T) T {
	if !is.Valid(value) {
		return d
	}
	return value
}

// FileExt extracts the file extension from the given file path.
// includes query parameters or fragments.
func FileExt(filePath string) string {
	ext := ""
	if len(filePath) > 0 {
		for i := len(filePath) - 1; i >= 0; i-- {
			if filePath[i] == '.' {
				ext = filePath[i+1:]
				break
			}
			if filePath[i] == '/' || filePath[i] == '\\' {
				break
			}
		}
	}
	return ext
}

// ValueOrIf returns the value v if the expression exp is true,
// otherwise returns the default value d.
func ValueOrIf[T any](exp bool, v T, d T) T {
	if exp {
		return v
	}
	return d
}

// ValIfOk returns the received pointer to value v if the expression exp is true,
// otherwise returns nil. If v is nil, it also returns nil.
func ValIfOk[T any](v *T, exp bool) *T {
	if exp {
		return v
	}
	return nil
}
