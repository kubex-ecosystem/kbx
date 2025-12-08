package kbx

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/google/uuid"
	gl "github.com/kubex-ecosystem/logz"
)

type IGlobalRef interface {
	GetID() uuid.UUID
	GetName() string
	SetName(name string)
	String() string
	GetGlobalRef() *GlobalRef
}

// GlobalRef is a struct that holds the GlobalRef ID and name.
type GlobalRef struct {
	// refID is the unique identifier for this context.
	ID uuid.UUID
	// refName is the name of the context.
	Name string
}

// newGlobalRef is a function that creates a new GlobalRef instance.
func newGlobalRef(name string) *GlobalRef {
	if name == "" {
		pc, _, line, ok := runtime.Caller(1)
		if ok {
			fn := runtime.FuncForPC(pc)
			name = fmt.Sprintf("%s:%d", fn.Name(), line)
		} else {
			name = "unknown"
		}
	}
	return &GlobalRef{
		ID:   uuid.New(),
		Name: name,
	}
}

// NewGlobalRef is a function that creates a new IGlobalRef instance.
func NewGlobalRef(name string) IGlobalRef {
	return newGlobalRef(name)
}

// String is a method that returns the string representation of the GlobalRef.
func (r *GlobalRef) String() string {
	return fmt.Sprintf("ID: %s, Name: %s", r.ID.String(), r.Name)
}

// GetID is a method that returns the ID of the GlobalRef.
func (r *GlobalRef) GetID() uuid.UUID {
	if r == nil {
		gl.Log("error", "GetID: GlobalRef does not exist (", reflect.TypeFor[GlobalRef]().String(), ")")
		return uuid.Nil
	}
	return r.ID
}

// GetName is a method that returns the name of the GlobalRef.
func (r *GlobalRef) GetName() string {
	if r == nil {
		gl.Log("error", "GetName: GlobalRef does not exist (", reflect.TypeFor[GlobalRef]().String(), ")")
		return ""
	}
	return r.Name
}

// SetName is a method that sets the name of the GlobalRef.
func (r *GlobalRef) SetName(name string) {
	if r == nil {
		gl.Log("error", "SetName: GlobalRef does not exist (", reflect.TypeFor[GlobalRef]().String(), ")")
		return
	}
	r.Name = name
}

// GetGlobalRef is a method that returns the GlobalRef struct (non-interface).
func (r *GlobalRef) GetGlobalRef() *GlobalRef {
	if r == nil {
		gl.Errorf("GetGlobalRef: GlobalRef does not exist (%s)", reflect.TypeFor[GlobalRef]().String())
		return nil
	}
	return r
}
