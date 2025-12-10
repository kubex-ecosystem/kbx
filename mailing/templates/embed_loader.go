package templates

import (
	"embed"
	"fmt"
)

// TemplateLoader define como obter o HTML de um template.
type TemplateLoader interface {
	LoadHTML(name string) (string, error)
}

// EmbedTemplateLoader carrega templates a partir de um embed.FS.
// Espera estrutura email/<template>/content.html.
type EmbedTemplateLoader struct {
	FS embed.FS
}

func (l *EmbedTemplateLoader) LoadHTML(name string) (string, error) {
	if l == nil {
		return "", fmt.Errorf("template loader is nil")
	}
	b, err := l.FS.ReadFile("email/" + name + "/content.html")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
