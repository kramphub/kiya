package backend

import (
	"context"
	"time"
)

type Backend interface {
	Get(ctx context.Context, p *Profile, key string) ([]byte, error)
	List(ctx context.Context, p *Profile) ([]Key, error)
	CheckExists(ctx context.Context, p *Profile, key string) (bool, error)
	Put(ctx context.Context, p *Profile, key, value string, overwrite bool) error
	Delete(ctx context.Context, p *Profile, key string) error
	SetParameter(key string, value interface{})
	Close() error
}

type Key struct {
	Name      string
	CreatedAt time.Time
	Owner     string
	Info      string
}

// Profile describes a single profile in a .kiya configuration
type Profile struct {
	Backend string
	// General
	Label               string
	Location            string
	SecretRunes         []rune
	AllowedCharacters   string `json:"allowedCharacters"`   // if set then this is the set of characters that will be used to generate secrets for this profile
	AutoCopyEnabled     bool   `json:"autoCopyEnabled"`     // if true then the secret of a single list result will be copied to clipboard
	PromptForSecretLine bool   `json:"promptForSecretLine"` // if true then you must enter a number to run the command on that line

	// GCP
	ProjectID string
	Keyring   string
	CryptoKey string
	Bucket    string

	// AWS
	AWSProfile string `json:"awsprofile"`

	// Vault
	VaultUrl       string
	VaultMountPath string
}
