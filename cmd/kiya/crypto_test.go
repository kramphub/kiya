package main

import "testing"

func TestGenerateSecret(t *testing.T) {
	secret := generateSecret()
	t.Logf("secret: %s", secret)
}
