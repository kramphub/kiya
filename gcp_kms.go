package kiya

import (
	"context"
	"encoding/base64"
	"fmt"

	cloudkms "cloud.google.com/go/kms/apiv1"
	"github.com/emicklei/tre"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

// GetDecryptedValue decrypts an encrypted value via Google KMS
func GetDecryptedValue(kmsService *cloudkms.KeyManagementClient, target Profile, cipherText string) (string, error) {
	path := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
		target.ProjectID,
		target.Location,
		target.Keyring,
		target.CryptoKey)

	decryptReq := &kmspb.DecryptRequest{
		Name:       path,
		Ciphertext: []byte(cipherText),
	}
	resp, err := kmsService.Decrypt(context.Background(), decryptReq)
	if err != nil {
		return "", tre.New(err, "failed to decrypt", "path", path)
	}
	data, err := base64.StdEncoding.DecodeString(string(resp.Plaintext))
	if err != nil {
		return "", tre.New(err, "failed to base64 decode")
	}
	return string(data), nil
}

// GetEncryptedValue converts a plain text to a Google KMS encrypted text
func GetEncryptedValue(kmsService *cloudkms.KeyManagementClient, target Profile, plainText string) (string, error) {
	path := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
		target.ProjectID,
		target.Location,
		target.Keyring,
		target.CryptoKey)
	encryptReq := &kmspb.EncryptRequest{
		Name: path,
		Plaintext: []byte(base64.RawURLEncoding.EncodeToString([]byte(plainText))),
	}
	resp, err := kmsService.Encrypt(context.Background(), encryptReq)
	if err != nil {
		return "", tre.New(err, "failed to encrypt")
	}
	return string(resp.Ciphertext), nil
}
