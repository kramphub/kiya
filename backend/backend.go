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
	Backend         string
	Label           string
	ProjectID       string
	Location        string
	Keyring         string
	CryptoKey       string
	Bucket          string
	VaultUrl        string
	SecretRunes     []rune
	AutoCopyEnabled bool // if true then the secret of a single list result will be copied to clipboard
}
