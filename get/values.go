package get

import (
	"encoding/json"
	"os"
	"reflect"

	"github.com/kubex-ecosystem/kbx/is"
)

func ValueOr[T any](value T, d T) (T, reflect.Type) {
	if !is.Valid(value) {
		return d, reflect.TypeFor[T]()
	}
	return value, reflect.TypeFor[T]()
}

func ValErrOr[T any](fn func() (T, error), d T) T {
	if is.Safe(fn, false) {
		value, err := fn()
		if err != nil || !is.Valid(value) {
			return d
		}
		return value
	}
	return d
}

func ValOrType[T any](value T, d T) T {
	if !is.Valid(value) {
		return d
	}
	return value
}

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
	if is.Valid(result) {
		return result
	}
	return result
}

func ValOrAny[T any](value T, d T) T {
	if !is.Valid(value) {
		return d
	}
	return value
}
