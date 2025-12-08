// Package descriptors provides various descriptor types used across the application.
package descriptors

import "embed"

//go:embed all:*.json
var manifestJSONFiles embed.FS

func GetManifestJSONFiles() embed.FS { return manifestJSONFiles }
