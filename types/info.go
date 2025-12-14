package types

import "os"

var (
	KubexManifest *ManifestImpl
)

type ManifestImpl struct {
	Name            string   `json:"name" yaml:"name" toml:"name" mapstructure:"name"`
	Application     string   `json:"application" yaml:"application" toml:"application" mapstructure:"application"`
	Version         string   `json:"version" yaml:"version" toml:"version" mapstructure:"version"`
	Private         bool     `json:"private" yaml:"private" toml:"private" mapstructure:"private"`
	Published       bool     `json:"published" yaml:"published" toml:"published" mapstructure:"published"`
	Aliases         []string `json:"aliases" yaml:"aliases" toml:"aliases" mapstructure:"aliases"`
	Repository      string   `json:"repository" yaml:"repository" toml:"repository" mapstructure:"repository"`
	Homepage        string   `json:"homepage" yaml:"homepage" toml:"homepage" mapstructure:"homepage"`
	Description     string   `json:"description" yaml:"description" toml:"description" mapstructure:"description"`
	GoVersion       string   `json:"go_version" yaml:"go_version" toml:"go_version" mapstructure:"go_version"`
	Main            string   `json:"main" yaml:"main" toml:"main" mapstructure:"main"`
	Bin             string   `json:"bin" yaml:"bin" toml:"bin" mapstructure:"bin"`
	Author          string   `json:"author" yaml:"author" toml:"author" mapstructure:"author"`
	Organization    string   `json:"organization" yaml:"organization" toml:"organization" mapstructure:"organization"`
	License         string   `json:"license" yaml:"license" toml:"license" mapstructure:"license"`
	Keywords        []string `json:"keywords" yaml:"keywords" toml:"keywords" mapstructure:"keywords"`
	Platforms       []string `json:"platforms" yaml:"platforms" toml:"platforms" mapstructure:"platforms"`
	Dependencies    []string `json:"dependencies" yaml:"dependencies" toml:"dependencies" mapstructure:"dependencies"`
	HealthcheckType string   `json:"healthcheck_type" yaml:"healthcheck_type" toml:"healthcheck_type" mapstructure:"healthcheck_type"`
	HealthcheckURL  string   `json:"healthcheck_url" yaml:"healthcheck_url" toml:"healthcheck_url" mapstructure:"healthcheck_url"`
	HealthcheckCmd  string   `json:"healthcheck_cmd" yaml:"healthcheck_cmd" toml:"healthcheck_cmd" mapstructure:"healthcheck_cmd"`
}

func (m *ManifestImpl) GetName() string        { return m.Name }
func (m *ManifestImpl) GetVersion() string     { return m.Version }
func (m *ManifestImpl) GetAliases() []string   { return m.Aliases }
func (m *ManifestImpl) GetRepository() string  { return m.Repository }
func (m *ManifestImpl) GetHomepage() string    { return m.Homepage }
func (m *ManifestImpl) GetDescription() string { return m.Description }
func (m *ManifestImpl) GetMain() string        { return m.Main }
func (m *ManifestImpl) GetBin() string {
	if m.Bin == "" {
		m.Bin, _ = os.Executable()
	}
	return m.Bin
}
func (m *ManifestImpl) GetAuthor() string      { return m.Author }
func (m *ManifestImpl) GetLicense() string     { return m.License }
func (m *ManifestImpl) GetKeywords() []string  { return m.Keywords }
func (m *ManifestImpl) GetPlatforms() []string { return m.Platforms }
func (m *ManifestImpl) IsPrivate() bool        { return m.Private }

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
