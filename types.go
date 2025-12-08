package kbx

import "reflect"

type ObjectsInfo struct {
	TypeName     string
	Kind         reflect.Kind
	IsParsable   bool
	IsNullable   bool
	IsAcceptable bool
	IsExecutable bool
}
