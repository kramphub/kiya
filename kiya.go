package kiya

import (
	"context"
	"time"

	cloudstore "cloud.google.com/go/storage"
	"github.com/emicklei/tre"
	"google.golang.org/api/cloudkms/v1"
	"google.golang.org/api/iterator"
)

func List(storageService *cloudstore.Client, target Profile) ([]Key, error) {
	ctx := context.Background()
	bucket := storageService.Bucket(target.Bucket)
	query := &cloudstore.Query{}
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
			Owner:     next.Owner,
		})
	}
	return keys, nil
}

func Move(kmsService *cloudkms.Service,
	storageService *cloudstore.Client,
	source Profile,
	sourceKey string,
	target Profile,
	targetKey string) error {

	// fetch value for key from source
	sourceValue, err := GetValueByKey(kmsService, storageService, sourceKey, source)
	if err != nil {
		return tre.New(err, "get source key failed", "key", sourceKey)
	}

	if err := PutSecret(kmsService, storageService, target, targetKey, sourceValue); err != nil {
		return tre.New(err, "save key failed", targetKey)
	}
	// delete key from source
	err = DeleteSecret(kmsService, storageService, source, sourceKey)
	return tre.New(err, "could not delete key", targetKey)
}

func CheckSecretExists(storageService *cloudstore.Client, target Profile, key string) bool {
	_, err := LoadSecret(storageService, target, key)
	if err == nil {
		return true
	}
	return false
}

// PutSecret encrypts the given value and stores it
func PutSecret(kmsService *cloudkms.Service,
	storageService *cloudstore.Client,
	target Profile,
	key,
	value string) error {

	encryptedValue, err := GetEncryptedValue(kmsService, target, value)
	if err != nil {
		return tre.New(err, "failed to fetch encrypted value", "key", key)
	}
	err = StoreSecret(storageService, target, key, encryptedValue)
	return tre.New(err, "store secret failed", "key", key, "encryptedValue", encryptedValue)
}

// DeleteSecret removes a key from the bucket
func DeleteSecret(kmsService *cloudkms.Service, storageService *cloudstore.Client, target Profile, key string) error {
	_, err := GetValueByKey(kmsService, storageService, key, target)
	if err != nil {
		return tre.New(err, "delete failed", "key", key, "err", err)
	}

	bucket := storageService.Bucket(target.Bucket)
	if _, err := bucket.Attrs(context.Background()); err != nil {
		return tre.New(err, "bucket does not exist", "bucket", target.Bucket)
	}
	err = bucket.Object(key).Delete(context.Background())
	return tre.New(err, "failed to delete secret", "key", key)
}


