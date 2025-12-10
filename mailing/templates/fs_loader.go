package templates

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileSystemTemplateLoader lê templates a partir de um diretório base.
// Estrutura esperada: <base>/<template>/content.html.
type FileSystemTemplateLoader struct {
	BasePath string
}

func (l *FileSystemTemplateLoader) LoadHTML(name string) (string, error) {
	if l == nil {
		return "", fmt.Errorf("template loader is nil")
	}
	p := filepath.Join(l.BasePath, name, "content.html")
	b, err := os.ReadFile(p)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
