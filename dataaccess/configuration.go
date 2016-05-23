package dataaccess

import "math/rand"

// Configuration retrieves the configuration for the application.
type Configuration struct {
	ID                   string `bson:"_id" json:"id"`
	SessionEncryptionKey []byte `json:"sessionEncryptionKey"`
	// SetSecureFlag sets whether cookies should be issued with the secure flag set.
	// When the secure flag is set, cookies cannot be transmitted over HTTP.
	// SSL must already be in place before this option is set.
	SetSecureFlag bool `json:"setSecureFlag"`
}

// NewConfiguration creates a new configuration file.
func NewConfiguration(sessionEncryptionKey []byte) *Configuration {
	return &Configuration{
		ID:                   "configuration",
		SessionEncryptionKey: sessionEncryptionKey,
	}
}

func createSessionEncryptionKey() []byte {
	key := make([]byte, 32)
	rand.Read(key)
	return key
}
