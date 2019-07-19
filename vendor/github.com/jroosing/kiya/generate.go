package main

import (
	"bytes"
	"crypto/rand"
	"math/big"
)

// default set contains characters that do not required URL encoding
// the kiya configuration can override this set per profile.
const defaultSecreteCharSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-~"

// GenerateSecret composes a random secrets using runes from a give set.
func GenerateSecret(length int, runes []rune) (string, error) {
	if len(runes) == 0 {
		runes = []rune(defaultSecreteCharSet)
	}
	var buffer bytes.Buffer
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(runes))))
		if err != nil {
			return "", err
		}
		buffer.WriteRune(runes[n.Int64()])
	}
	return buffer.String(), nil
}
