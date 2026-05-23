package interfaces

// IKeyringService defines the contract for secret storage backends.
type IKeyringService interface {
	StorePassword(password string) error
	RetrievePassword() (string, error)
}

// ICryptoService defines the contract for symmetric encryption operations.
type ICryptoService interface {
	GenerateKey() ([]byte, error)
	Encrypt(data, key []byte) (string, string, error)
	Decrypt(data, key []byte) ([]byte, string, error)
	IsEncrypted(data []byte) bool
	EncodeBase64(data []byte) string
	DecodeBase64(s string) ([]byte, error)
	IsBase64String(s string) bool
}
