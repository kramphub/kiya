package kiya

import (
	"fmt"
	"io/ioutil"

	cloudstore "cloud.google.com/go/storage"
	"github.com/emicklei/tre"
	"golang.org/x/net/context"
)

func storeSecret(storageService *cloudstore.Client, target Profile, key, encryptedValue string) error {
	bucket := storageService.Bucket(target.Bucket)
	if _, err := bucket.Attrs(context.Background()); err != nil {
		tre.New(err, "bucket does not exist", "bucket", target.Bucket)
	}
	w := bucket.Object(key).NewWriter(context.Background())
	defer w.Close()
	_, err := fmt.Fprintf(w, encryptedValue)
	return tre.New(err, "writing encrypted value failed", "encryptedValue", encryptedValue)
}

func deleteSecret(storageService *cloudstore.Client, target Profile, key string) error {
	bucket := storageService.Bucket(target.Bucket)
	if _, err := bucket.Attrs(context.Background()); err != nil {
		tre.New(err, "bucket does not exist", "bucket", target.Bucket)
	}
	if err := bucket.Object(key).Delete(context.Background()); err != nil {
		return tre.New(err, "failed to delete secret", "key", key)
	}
	return nil
}

func loadSecret(storageService *cloudstore.Client, target Profile, key string) (string, error) {
	bucket := storageService.Bucket(target.Bucket)
	r, err := bucket.Object(key).NewReader(context.Background())
	if err != nil {
		return "", tre.New(err, "failed to get bucket", "profile", target.Label, "key", key)
	}
	defer r.Close()
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return "", tre.New(err, "reading encrypted value failed", "profile", target.Label, "key", key)
	}
	return string(data), nil
}
