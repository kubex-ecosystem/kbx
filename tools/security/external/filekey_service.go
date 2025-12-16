// Package external implements a file-based, AES-GCM encrypted replacement for go-keyring.
package external

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"

	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/kubex-ecosystem/kbx"

	kbxGet "github.com/kubex-ecosystem/kbx/get"
	sci "github.com/kubex-ecosystem/kbx/tools/security/interfaces"
	gl "github.com/kubex-ecosystem/logz"
)

var (
	_ sci.IKeyringService = (*FileKeyringService)(nil)
)

// FileKeyringService is a drop-in replacement for KeyringService,
// maintaining the same contract and method signatures.
type FileKeyringService struct {
	keyringService kbx.GlobalRef
	keyringName    kbx.GlobalRef
	masterKey      []byte
	baseDir        string
}

// NewFileKeyringService creates a new encrypted file-based secret store.
func NewFileKeyringService(service, name string) sci.IKeyringService {
	return newFileKeyringService(service, name, nil)
}

func newFileKeyringService(service, name string, def any) *FileKeyringService {
	// m := kbxGet.EnvOr(strings.Join([]string{strings.ToUpper(types.KubexManifest.GetName()), "_MASTER_KEY"}, ""), "")
	dir := kbxGet.EnvOr("APP_SECRETS_DIR", "")//os.ExpandEnv(DefaultVaultDir))
	m := kbxGet.EnvOr("APP_MASTER_KEY", "")
	masterKeyPath := filepath.Join(dir, "master.key")

	// Tenta criar o diretório
	if err := os.MkdirAll(dir, 0o700); err != nil {
		gl.Log("warn", fmt.Sprintf("Failed to create secrets directory %s: %v, trying fallback", dir, err))
		dir = "/var/lib/canalize/secrets"
		if err := os.MkdirAll(dir, 0o700); err != nil {
			gl.Log("error", fmt.Sprintf("Failed to create fallback secrets directory %s: %v", dir, err))
		}
	}

	var raw []byte

	if m != "" {
		mk, err := os.Stat(masterKeyPath)
		if err == nil && mk.IsDir() {
			gl.Log("error", fmt.Sprintf("APP_MASTER_KEY points to a directory (%s), expected a file or base64 string", m))
		}
		if err == nil {
			data, err := os.ReadFile(masterKeyPath)
			if err != nil {
				gl.Log("error", fmt.Sprintf("failed to read APP_MASTER_KEY file (%s): %v", m, err))
			} else {
				m = strings.TrimSpace(string(data))
			}
		}

		raw, _ = base64.StdEncoding.DecodeString(m)
	} else {
		mk, err := os.ReadFile(masterKeyPath)
		mkStr := ""
		if err == nil {
			mkStr = strings.TrimSpace(string(mk))
		}
		mk = []byte(mkStr)

		if err != nil || len(mk) == 0 {
			// gerar novo master key
			newKey := make([]byte, 32)
			if _, err := rand.Read(newKey); err != nil {
				gl.Log("error", fmt.Sprintf("failed to generate new master key: %v", err))
			} else {
				mk = []byte(base64.StdEncoding.EncodeToString(newKey))
				err := os.WriteFile(masterKeyPath, []byte(mk), 0o600)
				if err != nil {
					gl.Log("error", fmt.Sprintf("failed to write new master key to file (%s): %v", masterKeyPath, err))
				} else {
					gl.Log("info", fmt.Sprintf("new master key generated and stored at %s", masterKeyPath))
				}
			}
		}

		// raw, _ := base64.StdEncoding.DecodeString(mk)
		dir := kbxGet.EnvOr("APP_SECRETS_DIR", "")
		if dir == "" {
			dir = "/var/lib/canalize/secrets"
		}
		_ = os.MkdirAll(dir, 0o700)
	}
	return &FileKeyringService{
		keyringService: kbx.NewGlobalRef(service),
		keyringName:    kbx.NewGlobalRef(name),
		masterKey:      raw,
		baseDir:        dir,
	}
}

func (k *FileKeyringService) StorePassword(password string) error {
	if password == "" {
		gl.Log("error", "key cannot be empty")
		return gl.Errorf("key cannot be empty")
	}
	enc, err := k.encrypt([]byte(password))
	if err != nil {
		return gl.Errorf("error encrypting password: %v", err)
	}
	path := filepath.Join(k.baseDir, fmt.Sprintf("%s_%s.secret", k.keyringService.GetName(), k.keyringName.GetName()))
	if err := os.WriteFile(path, []byte(enc), 0o600); err != nil {
		return gl.Errorf("error storing key: %v", err)
	}
	gl.Log("debug", fmt.Sprintf("key stored successfully: %s", k.keyringName.GetName()))
	return nil
}

func (k *FileKeyringService) RetrievePassword() (string, error) {
	path := filepath.Join(k.baseDir, fmt.Sprintf("%s_%s.secret", k.keyringService.GetName(), k.keyringName.GetName()))
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", os.ErrNotExist
		}
		gl.Log("debug", fmt.Sprintf("error reading key: %v", err))
		return "", gl.Errorf("error retrieving key: %v", err)
	}
	plain, err := k.decrypt(string(data))
	if err != nil {
		return "", gl.Errorf("error decrypting key: %v", err)
	}
	return string(plain), nil
}

// --- internal helpers ---

func (k *FileKeyringService) encrypt(plain []byte) (string, error) {
	hash := sha256.Sum256(k.masterKey)
	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return "", err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ct := aead.Seal(nil, nonce, plain, nil)
	full := append(nonce, ct...)
	return base64.StdEncoding.EncodeToString(full), nil
}

func (k *FileKeyringService) decrypt(ciphertext string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(k.masterKey)
	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(data) < aead.NonceSize() {
		return nil, errors.New("invalid ciphertext")
	}
	nonce, ct := data[:aead.NonceSize()], data[aead.NonceSize():]
	return aead.Open(nil, nonce, ct, nil)
}
