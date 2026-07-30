package providers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	kbxTypes "github.com/kubex-ecosystem/kbx/types"
)

// Este é o teste que importa de verdade: prova que a chave do usuário CHEGA no
// header da chamada HTTP ao provedor. Antes desta correção, o campo existia,
// era preenchido pelo chamador e morria no caminho — a requisição saía sempre
// com a chave do servidor (ou seja, na conta do dono do servidor).
func TestBYOK_ChaveDoUsuarioChegaNoHeader(t *testing.T) {
	casos := []struct {
		nome     string
		novo     func(baseURL string) (kbxTypes.ProviderExt, error)
		header   string
		esperado func(chave string) string
	}{
		{
			nome: "groq",
			novo: func(u string) (kbxTypes.ProviderExt, error) {
				return NewGroqProvider("groq", u, "chave-do-servidor", "modelo")
			},
			header:   "Authorization",
			esperado: func(k string) string { return "Bearer " + k },
		},
		{
			nome: "openai",
			novo: func(u string) (kbxTypes.ProviderExt, error) {
				return NewOpenAIProvider("openai", u, "chave-do-servidor", "modelo")
			},
			header:   "Authorization",
			esperado: func(k string) string { return "Bearer " + k },
		},
		{
			nome: "anthropic",
			novo: func(u string) (kbxTypes.ProviderExt, error) {
				return NewAnthropicProvider("anthropic", u, "chave-do-servidor", "modelo")
			},
			header:   "x-api-key",
			esperado: func(k string) string { return k },
		},
	}

	for _, c := range casos {
		t.Run(c.nome+" usa a chave do usuário quando ela vem", func(t *testing.T) {
			recebido := make(chan string, 1)
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				recebido <- r.Header.Get(c.header)
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("data: [DONE]\n\n"))
			}))
			defer srv.Close()

			p, err := c.novo(srv.URL)
			if err != nil {
				t.Fatalf("criar provider: %v", err)
			}

			ch, err := p.Chat(context.Background(), kbxTypes.ChatRequest{
				Provider: c.nome,
				Model:    "modelo",
				Messages: []kbxTypes.Message{{Role: "user", Content: "oi"}},
				Headers:  map[string]string{kbxTypes.HeaderBYOK: "chave-do-usuario"},
			})
			if err != nil {
				t.Fatalf("Chat: %v", err)
			}
			drenar(ch)

			select {
			case got := <-recebido:
				if quero := c.esperado("chave-do-usuario"); got != quero {
					t.Errorf("header %s = %q, quero %q", c.header, got, quero)
				}
			case <-time.After(3 * time.Second):
				t.Fatal("o provider não chamou o servidor")
			}
		})

		t.Run(c.nome+" mantém a chave do servidor quando não há BYOK", func(t *testing.T) {
			recebido := make(chan string, 1)
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				recebido <- r.Header.Get(c.header)
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("data: [DONE]\n\n"))
			}))
			defer srv.Close()

			p, err := c.novo(srv.URL)
			if err != nil {
				t.Fatalf("criar provider: %v", err)
			}

			ch, err := p.Chat(context.Background(), kbxTypes.ChatRequest{
				Provider: c.nome,
				Model:    "modelo",
				Messages: []kbxTypes.Message{{Role: "user", Content: "oi"}},
				// sem Headers: comportamento de sempre precisa continuar igual
			})
			if err != nil {
				t.Fatalf("Chat: %v", err)
			}
			drenar(ch)

			select {
			case got := <-recebido:
				if quero := c.esperado("chave-do-servidor"); got != quero {
					t.Errorf("header %s = %q, quero %q", c.header, got, quero)
				}
			case <-time.After(3 * time.Second):
				t.Fatal("o provider não chamou o servidor")
			}
		})
	}
}

// drenar consome o canal até fechar, para a goroutine do provider terminar.
func drenar(ch <-chan kbxTypes.ChatChunk) {
	limite := time.After(3 * time.Second)
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return
			}
		case <-limite:
			return
		}
	}
}
