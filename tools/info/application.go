// Package info provides functionality to read and parse the application manifest.
package info

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"

	gl "github.com/kubex-ecosystem/logz"
)

// manifestJSONData and controlJSONData are not embedded here: callers must
// provide the manifest/control as []byte to the parsing functions. Only the
// JSON Schema files are embedded below.

//go:embed manifest.schema.json
var manifestSchemaData []byte

//go:embed control.schema.json
var controlSchemaData []byte

var (
	// cachedManifest Manifest
	kbxInfoInstance *mmanifest
	// cachedControl  *Control
	kbxControlInstance *Control
)

// var application Manifest

type Reference struct {
	Name            string `json:"name"`
	ApplicationName string `json:"application"`
	Bin             string `json:"bin"`
	Version         string `json:"version"`
}

type mmanifest struct {
	Manifest
	Name            string   `json:"name"`
	ApplicationName string   `json:"application"`
	Bin             string   `json:"bin"`
	Version         string   `json:"version"`
	Repository      string   `json:"repository"`
	Aliases         []string `json:"aliases,omitempty"`
	Homepage        string   `json:"homepage,omitempty"`
	Description     string   `json:"description,omitempty"`
	Main            string   `json:"main,omitempty"`
	Author          string   `json:"author,omitempty"`
	License         string   `json:"license,omitempty"`
	Keywords        []string `json:"keywords,omitempty"`
	Platforms       []string `json:"platforms,omitempty"`
	LogLevel        string   `json:"log_level,omitempty"`
	Debug           bool     `json:"debug,omitempty"`
	ShowTrace       bool     `json:"show_trace,omitempty"`
	Private         bool     `json:"private,omitempty"`
}

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

func (m *mmanifest) GetName() string        { return m.Name }
func (m *mmanifest) GetVersion() string     { return m.Version }
func (m *mmanifest) GetAliases() []string   { return m.Aliases }
func (m *mmanifest) GetRepository() string  { return m.Repository }
func (m *mmanifest) GetHomepage() string    { return m.Homepage }
func (m *mmanifest) GetDescription() string { return m.Description }
func (m *mmanifest) GetMain() string        { return m.Main }
func (m *mmanifest) GetBin() string {
	if m.Bin == "" {
		m.Bin, _ = os.Executable()
	}
	return m.Bin
}
func (m *mmanifest) GetAuthor() string      { return m.Author }
func (m *mmanifest) GetLicense() string     { return m.License }
func (m *mmanifest) GetKeywords() []string  { return m.Keywords }
func (m *mmanifest) GetPlatforms() []string { return m.Platforms }
func (m *mmanifest) IsPrivate() bool        { return m.Private }

// lazy cache
var (
	cachedManifest Manifest
	cachedControl  *Control
)

// GetManifest lazy, sem init() com side-effects
// GetManifest valida bytes do manifesto com o schema e retorna a interface Manifest.
// O conteúdo do manifesto deve ser fornecido via []byte (não há embed do manifesto neste pacote).
func GetManifest(manifest []byte) (Manifest, error) {
	if cachedManifest != nil {
		return cachedManifest, nil
	}

	if len(manifest) == 0 {
		return nil, gl.Errorf("manifest.json: no data provided")
	}

	// validar estrutural simples conforme manifest.schema.json: garantir campos required
	if err := validateManifestBasic(manifest); err != nil {
		return nil, gl.Errorf("manifest.json: validation error: %v", err)
	}

	var m mmanifest
	if err := json.Unmarshal(manifest, &m); err != nil {
		return nil, gl.Errorf("manifest.json: %v", err)
	}
	cachedManifest = &m
	return &m, nil
}

// FS secOrder quiser permitir override por FS externo:
type FS interface {
	ReadFile(name string) ([]byte, error)
}

func LoadFromFS(fs FS) (Manifest, Control, error) {
	var m Manifest
	var c Control
	if b, err := fs.ReadFile("manifest.json"); err == nil {
		// validar e decodificar via GetManifest
		mm, err := GetManifest(b)
		if err != nil {
			return nil, Control{}, gl.Errorf("manifest.json: %v", err)
		}
		m = mm
	} else {
		return nil, Control{}, gl.Errorf("manifest.json: %v", err)
	}
	if b, err := fs.ReadFile("control.json"); err == nil {
		ci, err := GetControl(b)
		if err != nil {
			return nil, Control{}, gl.Errorf("control.json: %v", err)
		}
		// if the returned interface is actually *Control, unwrap to value
		if ccptr, ok := ci.(*Control); ok {
			c = *ccptr
		} else {
			// fallback: unmarshal into concrete
			if err := json.Unmarshal(b, &c); err != nil {
				return nil, Control{}, gl.Errorf("control.json: %v", err)
			}
		}
	} else {
		return nil, Control{}, gl.Errorf("control.json: %v", err)
	}
	return m, c, nil
}

// ControlInterface expõe um subconjunto do Control para uso externo sem expor o concreto.
type ControlInterface interface {
	GetName() string
	GetVersion() string
}

// GetControl valida bytes do control com o schema e retorna uma interface (ControlInterface).
// O conteúdo do control deve ser fornecido via []byte (não há embed do control neste pacote).
func GetControl(control []byte) (ControlInterface, error) {
	if cachedControl != nil {
		return cachedControl, nil
	}

	if len(control) == 0 {
		return nil, gl.Errorf("control.json: no data provided")
	}

	// validação básica conforme control.schema.json
	if err := validateControlBasic(control); err != nil {
		return nil, gl.Errorf("control.json: validation error: %v", err)
	}

	var c Control
	if err := json.Unmarshal(control, &c); err != nil {
		return nil, gl.Errorf("control.json: %v", err)
	}
	cachedControl = &c
	return cachedControl, nil
}

// validateManifestBasic faz validações simples baseadas no schema: garante campos required
func validateManifestBasic(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	required := []string{"name", "application", "version", "private", "published", "aliases", "repository", "homepage", "description", "go_version", "main", "bin", "author", "organization", "license", "keywords", "platforms", "dependencies", "healthcheck_type", "healthcheck_url", "healthcheck_cmd"}
	for _, k := range required {
		if _, ok := raw[k]; !ok {
			return errors.New("missing required field: " + k)
		}
	}
	return nil
}

// validateControlBasic faz validações simples do control conforme schema
func validateControlBasic(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	ctrl, ok := raw["control"].(map[string]interface{})
	if !ok {
		return errors.New("missing required object: control")
	}
	required := []string{"schema_version", "module", "ipc", "bitreg", "kv", "seq", "epoch_ns"}
	for _, k := range required {
		if _, ok := ctrl[k]; !ok {
			return errors.New("missing required control field: " + k)
		}
	}
	return nil
}

func CLIBannerStyle(banners, descriptionArg []string, _ bool) map[string]string {
	var description, banner string

	if descriptionArg != nil {
		if strings.Contains(strings.Join(os.Args[0:], ""), "-h") {
			description = descriptionArg[0]
		} else {
			description = descriptionArg[1]
		}
	} else {
		description = ""
	}

	if kbxInfoInstance.Manifest == nil {
		return map[string]string{"banner": banner, "description": description}
	}

	if kbxInfoInstance.Manifest.GetDescription() != "" {
		description += kbxInfoInstance.Manifest.GetDescription()
	}

	bannerRandLen := len(banners)
	bannerRandIndex := rand.Intn(bannerRandLen)
	banner = fmt.Sprintf(banners[bannerRandIndex], "\033[1;34m", kbxInfoInstance.Manifest.GetVersion(), "\033[0m")

	return map[string]string{"banner": banner, "description": description}
}
