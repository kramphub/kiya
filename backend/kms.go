package backend

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"

	"cloud.google.com/go/storage"
	"github.com/emicklei/tre"
	"google.golang.org/api/cloudkms/v1"
	"google.golang.org/api/iterator"
)

type KMS struct {
	kmsService    *cloudkms.Service
	storageClient *storage.Client
}

func NewKMS(kmsService *cloudkms.Service, storageClient *storage.Client) *KMS {
	return &KMS{
		kmsService:    kmsService,
		storageClient: storageClient,
	}
}

func (b *KMS) Get(ctx context.Context, p *Profile, key string) ([]byte, error) {
	encryptedValue, err := b.loadSecret(p, key)
	if err != nil {
		log.Fatal(tre.New(err, "get failed", "key", key))
		return nil, err
	}

	decryptedValue, err := b.getDecryptedValue(p, encryptedValue)
	if err != nil {
		log.Fatal(tre.New(err, "get failed", "cipherText", encryptedValue))
		return nil, err
	}

	return decryptedValue, nil
}

func (b *KMS) CheckExists(ctx context.Context, p *Profile, key string) (bool, error) {
	bucket := b.storageClient.Bucket(p.Bucket)
	r, err := bucket.Object(key).NewReader(ctx)
	if err != nil {
		return false, tre.New(err, "failed to get bucket", "profile", p.Label, "key", key)
	}
	defer r.Close()

	_, err = ioutil.ReadAll(r)
	if err != nil {
		return false, tre.New(err, "reading encrypted value failed", "profile", p.Label, "key", key)
	}

	return true, nil
}

func (b *KMS) Put(ctx context.Context, p *Profile, key, value string) error {
	encryptedValue, err := b.getEncryptedValue(p, value)
	if err != nil {
		return tre.New(err, "failed to fetch encrypted value", "key", key)
	}

	if err := b.storeSecret(p, key, encryptedValue); err != nil {
		return tre.New(err, "store secret failed", "key", key, "encryptedValue", encryptedValue)
	}

	return nil
}

func (b *KMS) Delete(ctx context.Context, p *Profile, key string) error {
	_, err := b.Get(ctx, p, key)
	if err != nil {
		return tre.New(err, "delete failed", "key", key, "err", err)
	}

	bucket := b.storageClient.Bucket(p.Bucket)
	if _, err := bucket.Attrs(ctx); err != nil {
		return tre.New(err, "bucket does not exist", "bucket", p.Bucket)
	}

	err = bucket.Object(key).Delete(ctx)
	return tre.New(err, "failed to delete secret", "key", key)
}

func (b *KMS) List(ctx context.Context, p *Profile) ([]Key, error) {
	bucket := b.storageClient.Bucket(p.Bucket)
	query := &storage.Query{}
	it := bucket.Objects(ctx, query)

	var keys []Key
	for {
		next, err := it.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return nil, tre.New(err, "list failed")
		}
		keys = append(keys, Key{
			Name:      next.Name,
			CreatedAt: next.Created,
			Info:      fmt.Sprintf("creator: %s", next.Owner),
			Owner:     next.Owner,
		})
	}

	return keys, nil
}

func (b *KMS) Close() error {
	return b.storageClient.Close()
}

///

func (b *KMS) loadSecret(p *Profile, key string) ([]byte, error) {
	bucket := b.storageClient.Bucket(p.Bucket)
	r, err := bucket.Object(key).NewReader(context.Background())
	if err != nil {
		return nil, tre.New(err, "failed to get bucket", "profile", p.Label, "key", key)
	}
	defer r.Close()

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, tre.New(err, "reading encrypted value failed", "profile", p.Label, "key", key)
	}

	return data, nil
}

func (b *KMS) getDecryptedValue(p *Profile, bytes []byte) ([]byte, error) {
	decryptReq := &cloudkms.DecryptRequest{
		Ciphertext: string(bytes),
	}

	path := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
		p.ProjectID,
		p.Location,
		p.Keyring,
		p.CryptoKey,
	)

	resp, err := b.kmsService.Projects.Locations.KeyRings.CryptoKeys.Decrypt(path, decryptReq).Do()
	if err != nil {
		return nil, tre.New(err, "failed to decrypt", "path", path)
	}

	data, err := base64.StdEncoding.DecodeString(resp.Plaintext)
	if err != nil {
		return nil, tre.New(err, "failed to base64 decode")
	}

	return data, nil
}

func (b *KMS) getEncryptedValue(p *Profile, plainText string) (string, error) {
	encryptReq := &cloudkms.EncryptRequest{
		Plaintext: base64.RawURLEncoding.EncodeToString([]byte(plainText)),
	}

	path := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
		p.ProjectID,
		p.Location,
		p.Keyring,
		p.CryptoKey,
	)
	resp, err := b.kmsService.Projects.Locations.KeyRings.CryptoKeys.Encrypt(path, encryptReq).Do()
	if err != nil {
		return "", tre.New(err, "failed to encrypt")
	}

	return resp.Ciphertext, nil
}

func (b *KMS) storeSecret(p *Profile, key, encryptedValue string) error {
	bucket := b.storageClient.Bucket(p.Bucket)
	if _, err := bucket.Attrs(context.Background()); err != nil {
		return tre.New(err, "bucket does not exist", "bucket", p.Bucket)
	}

	w := bucket.Object(key).NewWriter(context.Background())
	defer w.Close()

	_, err := fmt.Fprintf(w, encryptedValue)
	return tre.New(err, "writing encrypted value failed", "encryptedValue", encryptedValue)
}

func (b *KMS) SetParameter(key string, value interface{}) {
	// noop
}
