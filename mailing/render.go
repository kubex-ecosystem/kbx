package mailing

import (
	"bytes"
	"text/template"
)

// RenderHTML aplica data sobre um template HTML simples.
// Usa text/template para permitir placeholders {{ .Field }}.
func RenderHTML(tmpl string, data any) (string, error) {
	parsed, err := template.New("email").Parse(tmpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := parsed.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
