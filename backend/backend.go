package backend

import (
	"context"
	"time"
)

type Backend interface {
	Get(ctx context.Context, p *Profile, key, version string) ([]byte, error)
	List(ctx context.Context, p *Profile) ([]*Key, error)
	CheckExists(ctx context.Context, p *Profile, key string) (bool, error)
	Put(ctx context.Context, p *Profile, key, value string) error
	Delete(ctx context.Context, p *Profile, key string) error
}
