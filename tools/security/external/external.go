package external

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	kbx "github.com/kubex-ecosystem/kbx"
	sci "github.com/kubex-ecosystem/kbx/tools/security/interfaces"
)

// FileKeyringService is an encrypted file-based secret store implementing IKeyringService.
type FileKeyringService struct {
	keyringService kbx.GlobalRef
	keyringName    kbx.GlobalRef
	masterKey      []byte
	baseDir        string
}

var _ sci.IKeyringService = (*FileKeyringService)(nil)

// NewFileKeyringService creates a new encrypted file-based keyring service.
func NewFileKeyringService(service, name string) sci.IKeyringService {
	dir := os.ExpandEnv(os.Getenv("KUBEX_GNYX_SECRETS_DIR"))
	if dir == "" {
		dir = "/var/lib/kubex/secrets"
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		dir = os.TempDir()
	}

	masterKeyPath := filepath.Join(dir, "master.key")
	var raw []byte

	if data, err := os.ReadFile(masterKeyPath); err == nil {
		decoded, decErr := base64.StdEncoding.DecodeString(strings.TrimSpace(string(data)))
		if decErr == nil {
			raw = decoded
		}
	}

	if len(raw) == 0 {
		newKey := make([]byte, 32)
		if _, err := io.ReadFull(rand.Reader, newKey); err == nil {
			encoded := base64.StdEncoding.EncodeToString(newKey)
			_ = os.WriteFile(masterKeyPath, []byte(encoded), 0o600)
			raw = newKey
		}
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
		return fmt.Errorf("keyring: password cannot be empty")
	}
	enc, err := k.encrypt([]byte(password))
	if err != nil {
		return fmt.Errorf("keyring: encrypt: %w", err)
	}
	path := filepath.Join(k.baseDir, fmt.Sprintf("%s_%s.secret", k.keyringService.GetName(), k.keyringName.GetName()))
	return os.WriteFile(path, []byte(enc), 0o600)
}

func (k *FileKeyringService) RetrievePassword() (string, error) {
	path := filepath.Join(k.baseDir, fmt.Sprintf("%s_%s.secret", k.keyringService.GetName(), k.keyringName.GetName()))
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", os.ErrNotExist
		}
		return "", fmt.Errorf("keyring: read: %w", err)
	}
	plain, err := k.decrypt(string(data))
	if err != nil {
		return "", fmt.Errorf("keyring: decrypt: %w", err)
	}
	return string(plain), nil
}

func (k *FileKeyringService) encrypt(plain []byte) (string, error) {
	block, err := aes.NewCipher(pad32(k.masterKey))
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
	return base64.StdEncoding.EncodeToString(append(nonce, ct...)), nil
}

func (k *FileKeyringService) decrypt(ciphertext string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(pad32(k.masterKey))
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(data) < aead.NonceSize() {
		return nil, errors.New("keyring: invalid ciphertext")
	}
	return aead.Open(nil, data[:aead.NonceSize()], data[aead.NonceSize():], nil)
}

func pad32(key []byte) []byte {
	padded := make([]byte, 32)
	copy(padded, key)
	return padded
}
