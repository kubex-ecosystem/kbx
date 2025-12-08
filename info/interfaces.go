package info

import (
	"github.com/google/uuid"
)

type Manifest interface {
	GetName() string
	GetVersion() string
	GetAliases() []string
	GetRepository() string
	GetHomepage() string
	GetDescription() string
	GetMain() string
	GetBin() string
	GetAuthor() string
	GetLicense() string
	GetKeywords() []string
	GetPlatforms() []string
	IsPrivate() bool
}

type KubexModuleCLI interface {
	ID() uuid.UUID
	Module() string
	Alias() string
	ShortDescription() string
	LongDescription() string
	Usage() string
	Examples() []string
	Active() bool
	Execute() error
}
