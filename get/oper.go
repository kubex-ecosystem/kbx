package get

import (
	"github.com/kubex-ecosystem/kbx/tools"
)

// Mapper is an alias for tools.Mapper[T]
type Mapper[T any] = tools.Mapper[T]

// Loader creates a new Mapper[T] to load values from the specified source.
// Example:
//
//	type MyStruct struct {
//	    Name string
//	    Age  int
//	}
//	mapper := get.Loader[MyStruct]("env")
//
// Will create a Mapper[MyStruct] that
// can load MyStruct values from files and environment. Returns: *Mapper[MyStruct]
func Loader[T any](from string) *Mapper[T] { return tools.NewEmptyMapperType[T](from) }
