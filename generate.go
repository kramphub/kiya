package main

import (
	"math/big"
	"bytes"
	"crypto/rand"
)

func GenerateSecret(length int, runes []rune) (string, error) {
	if len(runes) == 0 {
		runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789~!@#$%^&*()_+`-={}|[]\\:\"<>?,./")
	}

	var buffer bytes.Buffer
	for i := 0; i < length; i ++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(runes))))
		if err != nil {
			return "", err
		}
		buffer.WriteRune(runes[n.Int64()])
	}

	return buffer.String(), nil
}
