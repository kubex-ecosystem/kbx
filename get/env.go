package get

import (
	"encoding/json"
	"os"
	"reflect"

	"github.com/kubex-ecosystem/kbx/is"
)

// EnvOr retrieves the value of the environment variable named by the key.
// If the variable is empty or not set, it returns the provided default value d.
func EnvOr(key, defaultValue string) string {
	if v := os.ExpandEnv(os.Getenv(key)); len(v) > 0 {
		return v
	}
	return defaultValue
}

// EnvOrType retrieves an environment variable by key and attempts to convert it to the specified type T.
// If the environment variable is not set or conversion fails, it returns the default value.
func EnvOrType[T any](key string, d T) T {
	value := os.ExpandEnv(os.Getenv(key))
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
