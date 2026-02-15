package interfaces

// ICryptoService defines the interface for cryptographic operations, including encryption, decryption, key generation, and encoding/decoding functionalities.
type ICryptoService interface {
	// Encrypt encrypts the given data using ChaCha20-Poly1305 algorithm
	// It ensures the data is decoded before encryption and the key is valid
	// Encrypt(data []byte, key []byte) (encryptedData string, nonce string, err error)
	Encrypt([]byte, []byte) (string, string, error)
	// Decrypt decrypts the given encrypted data using ChaCha20-Poly1305 algorithm
	// It ensures the data is decoded before decryption and the key is valid
	// Decrypt(encryptedData string, nonce string, key []byte) (decryptedData string, err error)
	Decrypt([]byte, []byte) (string, string, error)

	// GenerateKey generates a random key of default length (32 bytes for ChaCha20-Poly1305)
	GenerateKey() ([]byte, error)
	// GenerateKeyWithLength generates a random key of specified length in bytes
	GenerateKeyWithLength(int) ([]byte, error)

	// EncodeIfDecoded encodes a byte slice to Base64 URL encoding if it is not already encoded
	EncodeIfDecoded([]byte) (string, error)
	// DecodeIfEncoded decodes a Base64 URL encoded byte slice if it is encoded, otherwise returns the original byte slice
	DecodeIfEncoded([]byte) ([]byte, error)
	// EncodeBase64 encodes a byte slice to Base64 URL encoding
	EncodeBase64([]byte) string
	// DecodeBase64 decodes a Base64 URL encoded string to a byte slice
	DecodeBase64(string) ([]byte, error)

	// IsBase64String checks if a given string is a valid Base64 URL encoded string
	IsBase64String(string) bool
	// IsKeyValid checks if a given byte slice is a valid key for the encryption algorithm (e.g., 32 bytes for ChaCha20-Poly1305)
	IsKeyValid([]byte) bool
	// IsEncrypted checks if a given byte slice is likely to be encrypted data (e.g., by checking for Base64 encoding or specific patterns)
	IsEncrypted([]byte) bool
}
