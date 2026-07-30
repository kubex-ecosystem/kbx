package types

import "testing"

// ResolveKey é o ÚNICO ponto de decisão do BYOK. Se ele errar, todo provider
// erra junto — daí os casos serem exaustivos.
func TestResolveKey(t *testing.T) {
	const configurada = "chave-do-servidor"

	casos := []struct {
		nome    string
		headers map[string]string
		quero   string
	}{
		{"sem headers usa a configurada", nil, configurada},
		{"headers vazios usam a configurada", map[string]string{}, configurada},
		{"header vazio usa a configurada", map[string]string{HeaderBYOK: ""}, configurada},
		{"só espaços usa a configurada", map[string]string{HeaderBYOK: "   "}, configurada},
		{"outro header não confunde", map[string]string{"Authorization": "x"}, configurada},
		{"chave do usuário vence", map[string]string{HeaderBYOK: "chave-do-user"}, "chave-do-user"},
		{"espaços em volta são aparados", map[string]string{HeaderBYOK: "  k  "}, "k"},
	}

	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			got := ChatRequest{Headers: c.headers}.ResolveKey(configurada)
			if got != c.quero {
				t.Errorf("ResolveKey = %q, quero %q", got, c.quero)
			}
		})
	}
}

func TestUsesBYOK(t *testing.T) {
	semHeaders := ChatRequest{}
	emBranco := ChatRequest{Headers: map[string]string{HeaderBYOK: " "}}
	comChave := ChatRequest{Headers: map[string]string{HeaderBYOK: "k"}}

	if semHeaders.UsesBYOK() {
		t.Error("sem headers não é BYOK")
	}
	if emBranco.UsesBYOK() {
		t.Error("chave em branco não é BYOK")
	}
	if !comChave.UsesBYOK() {
		t.Error("com chave deveria ser BYOK")
	}
}
