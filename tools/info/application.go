// Package info provides functionality to read and parse the application manifest.
package info

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/kubex-ecosystem/kbx/load"
	"github.com/kubex-ecosystem/kbx/types"
	gl "github.com/kubex-ecosystem/logz"
	jsonschema "github.com/santhosh-tekuri/jsonschema/v5"
)

// manifestJSONData and controlJSONData are not embedded here: callers must
// provide the manifest/control as []byte to the parsing functions. Only the
// JSON Schema files are embedded below.

//go:embed manifest.schema.json
var manifestSchemaData []byte

var (
	// cachedManifest Manifest
	kbxInfoInstance Manifest
)

type MManifest = types.MManifest
type Manifest = types.Manifest

// GetManifest lazy, sem init() com side-effects
// GetManifest valida bytes do manifesto com o schema e retorna a interface Manifest.
// O conteúdo do manifesto deve ser fornecido via []byte (não há embed do manifesto neste pacote).
func GetManifest(manifest []byte, path string) (Manifest, error) {
	if kbxInfoInstance != nil {
		gl.Debug("Using cached manifest")
		return kbxInfoInstance, nil
	}

	var loaded *MManifest

	if len(manifest) == 0 && len(path) == 0 {
		return nil, gl.Errorf("manifest.json: no data provided")
	} else if len(manifest) == 0 && len(path) > 0 {
		// load from file
		data, err := load.LoadConfig[MManifest](path)
		if err != nil {
			return nil, gl.Errorf("manifest.json: %v", err)
		}
		loaded = &data
	} else {
		gl.Debug("Validating manifest.json with JSON Schema")
		// validar com JSON Schema completo
		comp := jsonschema.NewCompiler()
		if err := comp.AddResource("manifest.schema.json", bytes.NewReader(manifestSchemaData)); err != nil {
			gl.Notice("Failed to add manifest schema resource: %v", err)
			return nil, gl.Errorf("manifest.schema: %v", err)
		}
		sch, err := comp.Compile("manifest.schema.json")
		if err != nil {
			gl.Notice("Failed to compile manifest schema: %v", err)
			return nil, gl.Errorf("manifest.schema compile: %v", err)
		}
		if err := sch.Validate(bytes.NewReader(manifest)); err != nil {
			gl.Error("Manifest validation error: %v", err)
			return nil, gl.Errorf("manifest.json: validation error: %v", err)
		}

		var m = &MManifest{}
		if err := json.Unmarshal(manifest, &m); err != nil {
			return nil, gl.Errorf("manifest.json: %v", err)
		}
		loaded = m
	}

	kbxInfoInstance = loaded
	types.KubexManifest = loaded
	return loaded, nil
}

// FS secOrder quiser permitir override por FS externo:
type FS interface {
	ReadFile(name string) ([]byte, error)
}

func LoadFromFS(fs FS, path string) (Manifest, error) {
	var m *MManifest
	// var c Control
	if b, err := fs.ReadFile("manifest.json"); err == nil {
		// validar e decodificar via GetManifest
		mm, err := GetManifest(b, path)
		if err != nil {
			return nil, gl.Errorf("manifest.json: %v", err)
		}
		m = mm.(*MManifest)
	} else {
		return nil, gl.Errorf("manifest.json: %v", err)
	}
	return m, nil
}

// ControlInterface expõe um subconjunto do Control para uso externo sem expor o concreto.
type ControlInterface interface {
	GetName() string
	GetVersion() string
}

func CLIBannerStyle(banners, descriptionArg []string, _ bool) map[string]string {
	if kbxInfoInstance == nil {
		return map[string]string{"banner": "", "description": ""}
	}

	return GetDefinitions(banners, descriptionArg, false)
}

func GetDefinitions(banners, descriptionArg []string, hideBanner bool) map[string]string {
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

	if kbxInfoInstance.GetDescription() != "" {
		description += kbxInfoInstance.GetDescription()
	}

	bannerRandLen := len(banners)
	bannerRandIndex := rand.Intn(bannerRandLen)
	banner = fmt.Sprintf(banners[bannerRandIndex], "\033[1;34m", kbxInfoInstance.GetVersion(), "\033[0m")

	return map[string]string{"banner": banner, "description": description}
}
