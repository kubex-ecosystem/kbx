package storage

// ISecretStorage defines the contract for low-level secret persistence.
type ISecretStorage interface {
	StorePassword(password string) error
	RetrievePassword() (string, error)
}
