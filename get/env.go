package get

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"

	"github.com/kubex-ecosystem/kbx/is"
)

// EnvOr retrieves the value of the environment variable named by the key.
// If the variable is empty or not set, it returns the provided default value d.
func EnvOr(key, defaultValue string) string {
	if v := os.ExpandEnv(os.Getenv(key)); len(strings.TrimSpace(v)) > 0 {
		return v
	}
	return defaultValue
}

// ExpandIfEnv resolves a configuration value by checking for an environment variable
// if the value is an environment variable, it will be expanded
func ExpandIfEnv(value string) string {
	if envVal := EnvOr(value, value); !strings.EqualFold(value, envVal) {
		return envVal // Expanded EnvVar
	}
	return value // Hardcoded value
}

// EnvOrType retrieves an environment variable by key and attempts to convert it to the specified type T.
// If the environment variable is not set or conversion fails, it returns the default value.
func EnvOrType[T any](key string, obj T) T {
	envVal := EnvOr(key, "")
	if len(envVal) == 0 {
		return obj
	}
	if reflect.ValueOf(envVal).CanConvert(reflect.TypeFor[T]()) {
		return reflect.ValueOf(envVal).Convert(reflect.TypeFor[T]()).Interface().(T)
	}
	var result T
	if err := json.Unmarshal([]byte(envVal), &result); err != nil {
		return obj
	}
	if is.Safe(result, false) {
		return result
	}
	return result
}
