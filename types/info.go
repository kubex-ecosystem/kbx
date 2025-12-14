package types

import "os"

var (
	KubexManifest *MManifest
)

type MManifest struct {
	Manifest `json:"-" yaml:"-" toml:"-" mapstructure:"-"`

	Name            string   `json:"name,omitempty" yaml:"name,omitempty" toml:"name,omitempty" mapstructure:"name,omitempty"`
	Application     string   `json:"application,omitempty" yaml:"application,omitempty" toml:"application,omitempty" mapstructure:"application,omitempty"`
	Version         string   `json:"version,omitempty" yaml:"version,omitempty" toml:"version,omitempty" mapstructure:"version,omitempty"`
	Private         bool     `json:"private,omitempty" yaml:"private,omitempty" toml:"private,omitempty" mapstructure:"private,omitempty"`
	Published       bool     `json:"published,omitempty" yaml:"published,omitempty" toml:"published,omitempty" mapstructure:"published,omitempty"`
	Aliases         []string `json:"aliases,omitempty" yaml:"aliases,omitempty" toml:"aliases,omitempty" mapstructure:"aliases,omitempty"`
	Repository      string   `json:"repository,omitempty" yaml:"repository,omitempty" toml:"repository,omitempty" mapstructure:"repository,omitempty"`
	Homepage        string   `json:"homepage,omitempty" yaml:"homepage,omitempty" toml:"homepage,omitempty" mapstructure:"homepage,omitempty"`
	Description     string   `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty" mapstructure:"description,omitempty"`
	GoVersion       string   `json:"go_version,omitempty" yaml:"go_version,omitempty" toml:"go_version,omitempty" mapstructure:"go_version,omitempty"`
	Main            string   `json:"main,omitempty" yaml:"main,omitempty" toml:"main,omitempty" mapstructure:"main,omitempty"`
	Bin             string   `json:"bin,omitempty" yaml:"bin,omitempty" toml:"bin,omitempty" mapstructure:"bin,omitempty"`
	Author          string   `json:"author,omitempty" yaml:"author,omitempty" toml:"author,omitempty" mapstructure:"author,omitempty"`
	Organization    string   `json:"organization,omitempty" yaml:"organization,omitempty" toml:"organization,omitempty" mapstructure:"organization,omitempty"`
	License         string   `json:"license,omitempty" yaml:"license,omitempty" toml:"license,omitempty" mapstructure:"license,omitempty"`
	Keywords        []string `json:"keywords,omitempty" yaml:"keywords,omitempty" toml:"keywords,omitempty" mapstructure:"keywords,omitempty"`
	Platforms       []string `json:"platforms,omitempty" yaml:"platforms,omitempty" toml:"platforms,omitempty" mapstructure:"platforms,omitempty"`
	Dependencies    []string `json:"dependencies,omitempty" yaml:"dependencies,omitempty" toml:"dependencies,omitempty" mapstructure:"dependencies,omitempty"`
	HealthcheckType string   `json:"healthcheck_type,omitempty" yaml:"healthcheck_type,omitempty" toml:"healthcheck_type,omitempty" mapstructure:"healthcheck_type,omitempty"`
	HealthcheckURL  string   `json:"healthcheck_url,omitempty" yaml:"healthcheck_url,omitempty" toml:"healthcheck_url,omitempty" mapstructure:"healthcheck_url,omitempty"`
	HealthcheckCmd  string   `json:"healthcheck_cmd,omitempty" yaml:"healthcheck_cmd,omitempty" toml:"healthcheck_cmd,omitempty" mapstructure:"healthcheck_cmd,omitempty"`
}

func (m *MManifest) GetName() string        { return m.Name }
func (m *MManifest) GetVersion() string     { return m.Version }
func (m *MManifest) GetAliases() []string   { return m.Aliases }
func (m *MManifest) GetRepository() string  { return m.Repository }
func (m *MManifest) GetHomepage() string    { return m.Homepage }
func (m *MManifest) GetDescription() string { return m.Description }
func (m *MManifest) GetMain() string        { return m.Main }
func (m *MManifest) GetBin() string {
	if m.Bin == "" {
		m.Bin, _ = os.Executable()
	}
	return m.Bin
}
func (m *MManifest) GetAuthor() string      { return m.Author }
func (m *MManifest) GetLicense() string     { return m.License }
func (m *MManifest) GetKeywords() []string  { return m.Keywords }
func (m *MManifest) GetPlatforms() []string { return m.Platforms }
func (m *MManifest) IsPrivate() bool        { return m.Private }

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
