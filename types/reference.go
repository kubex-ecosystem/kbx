package types

import "github.com/google/uuid"

type GlobalRef struct {
	ID   uuid.UUID `json:"id,omitempty"`
	Name string    `json:"name,omitempty"`
}

func NewGlobalRef(name string) GlobalRef {
	return GlobalRef{
		ID:   uuid.New(),
		Name: name,
	}
}

func (gr *GlobalRef) GetGlobalRef() GlobalRef { return *gr }
func (gr *GlobalRef) GetName() string         { return gr.Name }
func (gr *GlobalRef) GetID() uuid.UUID        { return gr.ID }
func (gr *GlobalRef) SetName(name string)     { gr.Name = name }
func (gr *GlobalRef) SetID(id uuid.UUID)      { gr.ID = id }
func (gr *GlobalRef) String() string {
	return gr.Name + "-" + gr.ID.String()
}
