package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	ethyrSecCrp "github.com/kubex-ecosystem/ethyr/tools/security/crypto"
)

// CryptoService implements ICryptoService using AES-GCM.
type CryptoService struct {
	*ethyrSecCrp.CryptoService
}

// NewCryptoService returns a new CryptoService.
func NewCryptoService() *CryptoService {
	return &CryptoService{}
}

func (c *CryptoService) GenerateKey() ([]byte, error) {
	return ethyrSecCrp.NewCryptoService().GenerateKeyWithLength(32)
}

func (c *CryptoService) GenerateKeyWithLength(n int) ([]byte, error) {
	if n <= 0 {
		n = 32
	}
	key := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("crypto: generate key: %w", err)
	}
	return key, nil
}

func (c *CryptoService) Encrypt(data, key []byte) (string, string, error) {
	block, err := aes.NewCipher(key[:32])
	if err != nil {
		return "", "", err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", err
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", "", err
	}
	ct := aead.Seal(nonce, nonce, data, nil)
	encoded := base64.StdEncoding.EncodeToString(ct)
	return encoded, "", nil
}

func (c *CryptoService) Decrypt(data, key []byte) ([]byte, string, error) {
	block, err := aes.NewCipher(key[:32])
	if err != nil {
		return nil, "", err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, "", err
	}
	ns := aead.NonceSize()
	if len(data) < ns {
		return nil, "", fmt.Errorf("crypto: ciphertext too short")
	}
	plain, err := aead.Open(nil, data[:ns], data[ns:], nil)
	if err != nil {
		return nil, "", err
	}
	return plain, "", nil
}

func (c *CryptoService) IsEncrypted(data []byte) bool {
	return len(data) > 0 && data[0] == 0x00
}

func (c *CryptoService) EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(append([]byte{0x00}, data...))
}

func (c *CryptoService) DecodeBase64(s string) ([]byte, error) {
	raw, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	if len(raw) > 0 && raw[0] == 0x00 {
		return raw[1:], nil
	}
	return raw, nil
}

func (c *CryptoService) IsBase64String(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}
