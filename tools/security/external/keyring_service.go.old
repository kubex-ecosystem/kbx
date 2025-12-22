package external

import (
	"errors"
	"os"

	crp "github.com/kubex-ecosystem/kbx/tools/security/crypto"
	sci "github.com/kubex-ecosystem/kbx/tools/security/interfaces"
	gl "github.com/kubex-ecosystem/logz"
	"github.com/zalando/go-keyring"
)

var (
	kbxKeyringService *KeyringService
	kbxCryptoService  sci.ICryptoService
)

// KeyringService is a drop-in replacement for KeyringService,
// maintaining the same contract and method signatures.

type KeyringService struct {
	krgServiceName string
	krgName        string
	cryptoService  sci.ICryptoService
}

func newKeyringService(service, name string) *KeyringService {
	if kbxKeyringService == nil {
		kbxKeyringService = &KeyringService{}
	}
	if kbxCryptoService == nil {
		kbxCryptoService = crp.NewCryptoService()
	}
	kbxKeyringService = &KeyringService{
		krgServiceName: service,
		krgName:        name,
		cryptoService:  kbxCryptoService,
	}
	return kbxKeyringService
}
func NewKeyringService(service, name string) sci.IKeyringService {
	return newKeyringService(service, name)
}
func NewKeyringServiceType(service, name string) *KeyringService {
	return newKeyringService(service, name)
}

func (k *KeyringService) StorePassword(password string) error {
	if k == nil {
		gl.Log("fatal", "KeyringService is nil, trying to create a new one")
		return errors.New("KeyringService is nil")
	}
	if password == "" {
		gl.Log("error", "key cannot be empty")
		return gl.Errorf("key cannot be empty")
	}
	service := k.krgServiceName
	name := k.krgName
	gl.Debugf("storing key: service=%s, name=%s", service, name)
	if err := keyring.Set(service, name, password); err != nil {
		return gl.Errorf("error storing key: %v", err)
	}
	gl.Debugf("key stored successfully: %s", name)
	return nil
}
func (k *KeyringService) RetrievePassword() (string, error) {
	if k == nil {
		gl.Log("fatal", "KeyringService is nil, trying to create a new one")
		return "", errors.New("KeyringService is nil")
	}
	// if k.krgName == "" || k.krgServiceName == "" {
	// 	k.krgName = kbx.KeyringName
	// 	k.krgServiceName = kbx.KeyringService
	// }
	service := k.krgServiceName
	name := k.krgName
	if password, err := keyring.Get(service, name); err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			gl.Debugf("key not found: %s, %v", name, err)
			return "", os.ErrNotExist
		}
		gl.Debugf("error retrieving key: %v", err)
		return "", gl.Errorf("error retrieving key: %v", err)
	} else {
		gl.Debugf("key retrieved successfully: %s", name)
		return password, nil
	}
}
func (k *KeyringService) DeletePassword() error {
	if k == nil {
		gl.Log("fatal", "KeyringService is nil, trying to create a new one")
		return errors.New("KeyringService is nil")
	}
	service := k.krgServiceName
	name := k.krgName
	if err := keyring.Delete(service, name); err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			gl.Debugf("key not found for deletion: %s, %v", name, err)
			return os.ErrNotExist
		}
		gl.Debugf("error deleting key: %v", err)
		return gl.Errorf("error deleting key: %v", err)
	}
	gl.Debugf("key deleted successfully: %s", name)
	return nil
}
func (k *KeyringService) RetrieveOrCreatePassword() (string, error) {
	if k == nil {
		gl.Log("fatal", "KeyringService is nil, trying to create a new one")
		return "", errors.New("KeyringService is nil")
	}
	password, err := k.RetrievePassword()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) || len(password) == 0 {
			gl.Debugf("password not found, generating a new one")
			newPasswordBytes, err := k.cryptoService.GenerateKeyWithLength(32)
			if err != nil {
				return "", gl.Errorf("error generating new password: %v", err)
			}
			newPassword := k.cryptoService.EncodeBase64(newPasswordBytes)
			if err := k.StorePassword(newPassword); err != nil {
				return "", gl.Errorf("error storing new password: %v", err)
			}
			return newPassword, nil
		}
		return "", gl.Errorf("error retrieving password: %v", err)
	}
	return password, nil
}
