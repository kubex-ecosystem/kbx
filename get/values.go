// Package get provides utility functions for retrieving values with default fallbacks and type handling.
// It includes functions to get values from environment variables, handle errors, and extract file extensions.
package get

import (
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

// ValOrType returns the value if it's valid,
// otherwise returns the default value.
func ValOrType[T any](value T, d T) T {
	if !is.Valid(value) {
		return d
	}
	return value
}

// ValueOrCb returns the default value if valid,
// otherwise will execute the callback and validate its result
func ValueOrCb[T *any](value *T, fn func() (*T, error)) (*T, error) {
	if is.Valid(value) {
		return value, nil
	}
	return ValueOrIf(is.Safe(fn, false), fn, func() (*T, error) {
		return nil, gl.Errorf("")
	})()
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

// ValueOrIf returns the value v if the expression exp is true,
// otherwise returns the default value d.
func ValueOrIf[T any](exp bool, v T, d T) T {
	if exp {
		return v
	}
	return d
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

// ValIfOk returns the received pointer to value v if the expression exp is true,
// otherwise returns nil. If v is nil, it also returns nil.
func ValIfOk[T any](v *T, exp bool) *T {
	if exp {
		return v
	}
	return nil
}
