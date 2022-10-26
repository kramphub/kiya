package main

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/kramphub/kiya/backend"
)

// Backup is a backup of all keys in store.
type Backup struct {
	//Encrypted secret with public key and encoded as base64 string
	Secret    string `json:"secret"`
	Encrypted bool   `json:"encrypted"`
	Data      []byte `json:"data"`
}

// String returns a base64 String representation of the Backup.
func (b *Backup) String() string {
	buf := encodeToJson(b)
	return base64.URLEncoding.EncodeToString(buf)
}

// FromString returns a Backup from a string representation.
func (b *Backup) FromString(str string) {
	buf, err := base64.URLEncoding.DecodeString(str)
	if err != nil {
		log.Fatalf("[FATAL] decode from string failed: %s", err.Error())
	}

	err = json.Unmarshal(buf, b)
	if err != nil {
		log.Fatalf("[FATAL] decode JSON string failed: %s", err.Error())
	}
}

// SecretAsBytes returns the secret as bytes.
func (b *Backup) SecretAsBytes() []byte {
	buf, err := base64.URLEncoding.DecodeString(b.Secret)
	if err != nil {
		log.Fatalf("[FATAL] decode secret base64 string failed: %s", err.Error())
	}

	return buf
}

// commandBackup creates a backup of all keys in store.
func commandBackup(ctx context.Context, b backend.Backend, target backend.Profile, filter string) (*Backup, error) {
	items, err := getItems(ctx, b, target, filter)
	if err != nil {
		return nil, err
	}

	buf := encodeToJson(items)

	return &Backup{Data: buf}, nil
}

// getItems returns all keys in store.
func getItems(ctx context.Context, b backend.Backend, target backend.Profile, filter string) (map[string][]byte, error) {
	items := make(map[string][]byte)

	keys := commandList(ctx, b, &target, filter)
	totalKeys := len(keys)

	for i, key := range keys {
		buf, err := b.Get(ctx, &target, key.Name)
		if err != nil {
			fmt.Printf("error: get key '%s' failed, %s", key.Name, err.Error())
			continue
		}

		items[key.Name] = buf
		fmt.Printf("\rSaved keys: %d/%d", i+1, totalKeys)
	}
	fmt.Println()

	return items, nil
}

// getPublicKey returns the public key from file or store.
func getPublicKey(ctx context.Context, b backend.Backend, target backend.Profile, location, key string) (*rsa.PublicKey, error) {
	switch location {
	case "store":
		buf, err := b.Get(ctx, &target, key)
		if err != nil {
			return nil, fmt.Errorf("get public key '%s' failed, %w", key, err)
		}

		return exportPublicKeyFromPEMString(buf), nil
	case "file":
		fallthrough
	default:
		buf, err := os.ReadFile(key)
		if err != nil {
			return nil, fmt.Errorf("read public file '%s' failed, %w", key, err)
		}

		return exportPublicKeyFromPEMString(buf), nil
	}
}

// getPrivateKey returns the private key from file.
func getPrivateKey(ctx context.Context, b backend.Backend, target backend.Profile, location, key string) (*rsa.PrivateKey, error) {
	switch location {
	case "store":
		buf, err := b.Get(ctx, &target, key)
		if err != nil {
			return nil, fmt.Errorf("get private key '%s' failed, %w", key, err)
		}

		return exportPrivateKeyFromPEMString(buf), nil
	case "file":
		fallthrough
	default:
		buf, err := os.ReadFile(key)
		if err != nil {
			return nil, fmt.Errorf("read private file '%s' failed, %w", key, err)
		}

		return exportPrivateKeyFromPEMString(buf), nil
	}
}
