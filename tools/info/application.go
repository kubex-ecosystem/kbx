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

//go:embed control.schema.json
var controlSchemaData []byte

var (
	// cachedManifest Manifest
	kbxInfoInstance *ManifestImpl
	// cachedControl  *Control
	kbxControlInstance *Control
)

type Manifest = types.Manifest

type ManifestImpl struct {
	*types.ManifestImpl
}


// lazy cache
var (
	cachedManifest Manifest
	cachedControl  *Control
)

// GetManifest lazy, sem init() com side-effects
// GetManifest valida bytes do manifesto com o schema e retorna a interface Manifest.
// O conteúdo do manifesto deve ser fornecido via []byte (não há embed do manifesto neste pacote).
func GetManifest(manifest []byte, path string) (Manifest, error) {
	if kbxInfoInstance != nil {
		gl.Debug("Using cached manifest")
		return kbxInfoInstance, nil
	}

	if len(manifest) == 0 && len(path) == 0 {
		return nil, gl.Errorf("manifest.json: no data provided")
	} else if len(manifest) == 0 && len(path) > 0 {
		// load from file
		data, err := load.LoadConfig[ManifestImpl](path)
		if err != nil {
			return nil, gl.Errorf("manifest.json: %v", err)
		}
		if data == nil {
			return nil, gl.Errorf("manifest.json: no data loaded from %s", path)
		}
		return data, nil
	}

	// validar com JSON Schema completo
	comp := jsonschema.NewCompiler()
	if err := comp.AddResource("manifest.schema.json", bytes.NewReader(manifestSchemaData)); err != nil {
		return nil, gl.Errorf("manifest.schema: %v", err)
	}
	sch, err := comp.Compile("manifest.schema.json")
	if err != nil {
		return nil, gl.Errorf("manifest.schema compile: %v", err)
	}
	if err := sch.Validate(bytes.NewReader(manifest)); err != nil {
		return nil, gl.Errorf("manifest.json: validation error: %v", err)
	}

	var m ManifestImpl
	if err := json.Unmarshal(manifest, &m); err != nil {
		return nil, gl.Errorf("manifest.json: %v", err)
	}
	kbxInfoInstance = &m
	return &m, nil
}

// FS secOrder quiser permitir override por FS externo:
type FS interface {
	ReadFile(name string) ([]byte, error)
}

func LoadFromFS(fs FS, path string) (Manifest, Control, error) {
	var m Manifest
	var c Control
	if b, err := fs.ReadFile("manifest.json"); err == nil {
		// validar e decodificar via GetManifest
		mm, err := GetManifest(b, path)
		if err != nil {
			return nil, Control{}, gl.Errorf("manifest.json: %v", err)
		}
		m = mm
	} else {
		return nil, Control{}, gl.Errorf("manifest.json: %v", err)
	}
	// if b, err := fs.ReadFile("control.json"); err == nil {
	// 	ci, err := GetControl(b)
	// 	if err != nil {
	// 		return nil, Control{}, gl.Errorf("control.json: %v", err)
	// 	}
	// 	// if the returned interface is actually *Control, unwrap to value
	// 	if ccptr, ok := ci.(*Control); ok {
	// 		c = *ccptr
	// 	} else {
	// 		// fallback: unmarshal into concrete
	// 		if err := json.Unmarshal(b, &c); err != nil {
	// 			return nil, Control{}, gl.Errorf("control.json: %v", err)
	// 		}
	// 	}
	// } else {
	// 	return nil, Control{}, gl.Errorf("control.json: %v", err)
	// }
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
	// if cachedControl != nil {
	// 	return cachedControl, nil
	// }

	// if len(control) == 0 {
	// 	return nil, gl.Errorf("control.json: no data provided")
	// }

	// // validar com JSON Schema completo
	// comp := jsonschema.NewCompiler()
	// if err := comp.AddResource("control.schema.json", bytes.NewReader(controlSchemaData)); err != nil {
	// 	return nil, gl.Errorf("control.schema: %v", err)
	// }
	// sch, err := comp.Compile("control.schema.json")
	// if err != nil {
	// 	return nil, gl.Errorf("control.schema compile: %v", err)
	// }
	// if err := sch.Validate(bytes.NewReader(control)); err != nil {
	// 	return nil, gl.Errorf("control.json: validation error: %v", err)
	// }

	// var c Control
	// if err := json.Unmarshal(control, &c); err != nil {
	// 	return nil, gl.Errorf("control.json: %v", err)
	// }
	// cachedControl = &c
	// return cachedControl, nil
	gl.Alert("GetControl is not implemented yet")
	return nil, nil
}

func CLIBannerStyle(banners, descriptionArg []string, _ bool) map[string]string {
	if kbxInfoInstance == nil {
		return map[string]string{"banner": "", "description": ""}
	}

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
