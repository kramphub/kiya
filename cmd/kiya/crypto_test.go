package main

import (
	"bytes"
	"encoding/base64"
	"log"
	"testing"
)

func TestEncryptSecret(t *testing.T) {
	secret := generateSecret()

	priv, pub, err := generateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	encryptedSecret, err := encryptSecret(secret, pub)
	if err != nil {
		t.Fatal(err)
	}

	decryptedSecret, err := decryptSecret(encryptedSecret, priv)
	if err != nil {
		t.Fatal(err)
	}

	originalSecret, err := base64.URLEncoding.DecodeString(secret)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(originalSecret, decryptedSecret) {
		t.Fatalf("original secret and decrypted secret is not equla 😮")
	}
}

func TestEncryptFile(t *testing.T) {
	bak := Backup{
		Secret: generateSecret(),
		Data:   []byte{0x0, 0xf, 0xff, 0x54, 0x43},
	}

	buf, err := encrypt(bak.Data, bak.SecretAsBytes())
	if err != nil {
		t.Fatal(err)
	}

	buf, err = decrypt(buf, bak.SecretAsBytes())
	if err != nil {
		t.Fatal(err)
	}

	if len(bak.Data) != len(buf) {
		t.Fatalf("origin and decrypted data length is not equal")
	}

	if bak.Data[0] != buf[0] || bak.Data[1] != buf[1] || bak.Data[2] != buf[2] || bak.Data[3] != buf[3] || bak.Data[4] != buf[4] {
		t.Fatalf("data in origin and decrypted is not equals")
	}
}

func TestEncryptSecretAndFile(t *testing.T) {
	bak := Backup{
		Secret: generateSecret(),
		Data:   []byte("my super secret text"),
	}

	encryptedDataBuf, err := encrypt(bak.Data, bak.SecretAsBytes())
	if err != nil {
		t.Fatal(err)
	}

	priv, pub, err := generateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	encryptedSecret, err := encryptSecret(bak.Secret, pub)
	if err != nil {
		t.Fatal(err)
	}

	bakStr := (&Backup{Secret: encryptedSecret, Data: encryptedDataBuf}).String()

	bak2 := Backup{}
	bak2.FromString(bakStr)

	if encryptedSecret != bak2.Secret {
		log.Fatalf("secret is not equal")
	}

	decryptedSecret, err := decryptSecret(bak2.Secret, priv)
	if err != nil {
		t.Fatal(err)
	}

	dBuf, err := decrypt(bak2.Data, decryptedSecret)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(bak.Data, dBuf) {
		t.Fatalf("origin and decrypted data is not equal")
	}
}
