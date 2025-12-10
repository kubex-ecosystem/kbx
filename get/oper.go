package get

import (
	"github.com/kubex-ecosystem/kbx/tools"
)

type Mapper[T any] = tools.Mapper[T]

func Loader[T any](from string) *Mapper[T] { return tools.NewEmptyMapperType[T](from) }
