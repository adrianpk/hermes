package am

// contextKey is a private type used for context keys to avoid collisions.
type contextKey string

// Define context keys
const (
	EncryptionKeyCtxKey contextKey = "encryptionKey"
)
