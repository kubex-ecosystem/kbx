// Package info provides functionality to read and parse the application manifest.
package info

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/kubex-ecosystem/kbx/info/ctl"
	"github.com/kubex-ecosystem/kbx/info/descriptors"

	gl "github.com/kubex-ecosystem/logz"
)

var banners = []string{
	`
     |  /        |                 
     | /  |   |  __ \    _ \  \  / 
     . \  |   |  |   |   __/    <  
    _|\_\\__._| _.__/  \___| _/\_\ 
    %s%b - %s%s
`,
}

func GetDescriptions(descriptionArg []string, _ bool) map[string]string {
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

	manifest, err := GetManifest()
	if err != nil {
		description += ""
	} else {
		if manifest.GetDescription() != "" {
			description += manifest.GetDescription()
		}
	}

	bannerRandLen := len(banners)
	bannerRandIndex := rand.Intn(bannerRandLen)
	banner = fmt.Sprintf(banners[bannerRandIndex], "\033[1;34m", manifest.GetName(), manifest.GetVersion(), "\033[0m")

	return map[string]string{"banner": banner, "description": description}
}

// manifestJSONFile is the name of the manifest JSON file
const manifestJSONFile = "manifest.json"

// Embeds do arquivo manifest.json gerado na build.
var manifestJSONData []byte

func init() {
	var err error
	manifestJSONData, err = descriptors.GetManifestJSONFiles().ReadFile(manifestJSONFile)
	if err != nil {
		gl.Log("error", fmt.Sprintf("Failed to read embedded manifest.json: %v", err))
		manifestJSONData = []byte{}
	}
}

// var application Manifest

type Reference struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	ApplicationName string    `json:"application"`
	Bin             string    `json:"bin"`
	Version         string    `json:"version"`
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
	cachedControl  *ctl.Control
)

// GetManifest lazy, sem init() com side-effects
func GetManifest() (Manifest, error) {
	if cachedManifest != nil {
		return cachedManifest, nil
	}

	if len(manifestJSONData) == 0 {
		return nil, gl.Errorf("manifest.json: embed is empty")
	}

	var m mmanifest
	if err := json.Unmarshal(manifestJSONData, &m); err != nil {
		return nil, gl.Errorf("manifest.json: %v", err)
	}
	cachedManifest = &m
	return &m, nil
}

// FS secOrder quiser permitir override por FS externo:
type FS interface {
	ReadFile(name string) ([]byte, error)
}

func LoadFromFS(fs FS) (Manifest, ctl.Control, error) {
	var m Manifest
	var c ctl.Control
	if b, err := fs.ReadFile("manifest.json"); err == nil {
		if err := json.Unmarshal(b, &m); err != nil {
			return nil, ctl.Control{}, gl.Errorf("manifest.json: %v", err)
		}
	} else {
		return nil, ctl.Control{}, gl.Errorf("manifest.json: %v", err)
	}
	if b, err := fs.ReadFile("control.json"); err == nil {
		if err := json.Unmarshal(b, &c); err != nil {
			return nil, ctl.Control{}, gl.Errorf("control.json: %v", err)
		}
	} else {
		return nil, ctl.Control{}, gl.Errorf("control.json: %v", err)
	}
	return m, c, nil
}

// func GetControl() (*Control, error) {
// 	if cachedControl != nil {
// 		return cachedControl, nil
// 	}
// 	var c Control
// 	if len(controlJSONData) == 0 {
// 		return nil, gl.Errorf("control.json: embed is empty")
// 	}
// 	if err := json.Unmarshal(controlJSONData, &c); err != nil {
// 		return nil, gl.Errorf("control.json: %v", err)
// 	}
// 	cachedControl = &c
// 	return &c, nil
// }
