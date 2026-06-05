package storage

import (
	"fmt"

	vault "github.com/hashicorp/vault/api"

	gl "github.com/kubex-ecosystem/logz"
)

type VaultSecretStorage struct {
	Client *vault.Client
	Mount  string
	Path   string
	Field  string
}

func NewVaultSecretStorage(addr, token, mount, path, field string) (*VaultSecretStorage, error) {
	cfg := vault.DefaultConfig()
	cfg.Address = addr
	client, err := vault.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	client.SetToken(token)
	return &VaultSecretStorage{
		Client: client,
		Mount:  mount,
		Path:   path,
		Field:  field,
	}, nil
}

func (v *VaultSecretStorage) StorePassword(password string) error {
	data := map[string]interface{}{"data": map[string]interface{}{v.Field: password}}
	_, err := v.Client.Logical().Write(fmt.Sprintf("%s/data/%s", v.Mount, v.Path), data)
	return err
}

func (v *VaultSecretStorage) RetrievePassword() (string, error) { // pragma: allowlist secret
	secret, err := v.Client.Logical().Read(fmt.Sprintf("%s/data/%s", v.Mount, v.Path))
	if err != nil {
		return "", err
	}
	if secret == nil || secret.Data["data"] == nil { // pragma: allowlist secret
		return "", gl.Errorf("secret not found")
	}
	data := secret.Data["data"].(map[string]interface{}) // pragma: allowlist secret
	val, ok := data[v.Field].(string)                    // pragma: allowlist secret
	if !ok {
		return "", gl.Errorf("invalid field type")
	}
	return val, nil
}
