package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"os"
)

// generateSecret returns generated secret as base64 string
func generateSecret() string {
	key := make([]byte, 64)

	_, err := rand.Read(key)
	if err != nil {
		// handle error here
	}

	return base64.URLEncoding.EncodeToString(key)
}

func encryptFile(data []byte, secret []byte) ([]byte, error) {
	key := sha256.Sum256(secret)
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, fmt.Errorf("create new Cipher failed, %s", err.Error())
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return ciphertext, nil
}

func decryptFile(data []byte, secret []byte) ([]byte, error) {
	key := sha256.Sum256(secret)
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, fmt.Errorf("create new Cipher failed, %s", err.Error())
	}

	if len(data) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)

	return data, nil
}

func encryptSecret(secret string, publicKey *rsa.PublicKey) (encryptedSecret string, err error) {
	buf, err := base64.URLEncoding.DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("decode base64 string failed, %w", err)
	}

	buf, err = rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, buf, nil)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(buf), nil
}

func decryptSecret(secret string, privateKey *rsa.PrivateKey) ([]byte, error) {
	secretBytes, err := base64.URLEncoding.DecodeString(secret)
	if err != nil {
		return nil, fmt.Errorf("decode base64 string failed, %w", err)
	}
	buf, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, secretBytes, nil)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func generateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, fmt.Errorf("generate RSA key pair failed, %w", err)
	}
	return privateKey, &privateKey.PublicKey, nil
}

// exportPublicKeyAsPEM returns public key as a string in PEM format
func exportPublicKeyAsPEM(key *rsa.PublicKey) string {
	pemStr := string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(key),
		},
	))
	return pemStr
}

// exportPrivateKeyAsPEM returns private key as a string in PEM format
func exportPrivateKeyAsPEM(key *rsa.PrivateKey) string {
	pemStr := string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	))
	return pemStr

}

func exportPrivateKeyFromPEMString(pemStr []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(pemStr)
	key, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	return key
}

func exportPublicKeyFromPEMString(pemStr []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pemStr)
	key, _ := x509.ParsePKCS1PublicKey(block.Bytes)
	return key
}

func saveKeyToFile(keyPem, filename string) error {
	pemBytes := []byte(keyPem)
	return os.WriteFile(filename, pemBytes, 0400)
}
