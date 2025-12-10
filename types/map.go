// Package types provides type-related utilities and definitions.
package types

import "reflect"

var (
	KindMap = map[reflect.Kind]bool{
		reflect.Struct:    true,
		reflect.Map:       true,
		reflect.Slice:     true,
		reflect.Array:     true,
		reflect.Chan:      true,
		reflect.Interface: true,
		reflect.Ptr:       true,
		reflect.String:    true,
		reflect.Int:       true,
		reflect.Float32:   true,
		reflect.Float64:   true,
		reflect.Bool:      true,
		reflect.Uint:      true,
		reflect.Uint8:     true,
		reflect.Uint16:    true,
		reflect.Uint32:    true,
		reflect.Uint64:    true,
		reflect.Int8:      true,
		reflect.Int16:     true,
		reflect.Int32:     true,
		reflect.Int64:     true,
	}
)
