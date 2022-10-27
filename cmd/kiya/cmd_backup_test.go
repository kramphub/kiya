package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
	"testing/fstest"
)

func TestBackup_String(t *testing.T) {
	bak := Backup{
		Secret: "foo bar",
		Data:   []byte{0x0, 0x1f, 0x2, 0x1e},
	}

	bakStr := bak.String()

	bak2 := Backup{}
	bak2.FromString(bakStr)

	if bak.Secret != bak2.Secret {
		t.Fail()
	}

	if len(bak.Data) != len(bak2.Data) {
		t.Fail()
	}

	if bak.Data[0] != bak2.Data[0] || bak.Data[1] != bak2.Data[1] || bak.Data[2] != bak2.Data[2] || bak.Data[3] != bak2.Data[3] {
		t.Fail()
	}
}

func TestBackupWithoutEncryption(t *testing.T) {
	input := map[string]interface{}{
		"bar": "bar string",
	}

	buf, err := json.Marshal(input)

	require.NoError(t, err)

	backup := Backup{
		Secret:    "",
		Encrypted: false,
		Data:      buf,
	}

	mockFS := fstest.MapFS{
		"backup_test": {
			Data: []byte(backup.String()),
		},
	}

	backBuf, err := mockFS.ReadFile("backup_test")
	require.NoError(t, err)

	backup2 := Backup{}
	backup2.FromString(string(backBuf))

	require.Equal(t, len(backup.Data), len(backup2.Data))

	backupData := make(map[string]interface{})

	err = json.Unmarshal(backup2.Data, &backupData)
	require.NoError(t, err)

	bar := input["bar"]
	bar2 := backupData["bar"]

	require.Equal(t, bar, bar2)
}

func TestBackupWithEncryption(t *testing.T) {
	input, buf := setupTestData(t)

	secret := generateSecret()
	privateKey, publicKey, err := generateKeyPair()
	require.NoError(t, err)

	backup := Backup{
		Secret:    secret,
		Encrypted: true,
		Data:      buf,
	}

	encryptedBuf, err := encrypt(buf, backup.SecretAsBytes())
	require.NoError(t, err)
	backup.Data = encryptedBuf
	encryptedSecret, err := encryptSecret(secret, publicKey)
	require.NoError(t, err, "secret encryption failed")
	backup.Secret = encryptedSecret

	mockFS := fstest.MapFS{
		"backup_test": {
			Data: []byte(backup.String()),
		},
	}

	backBuf, err := mockFS.ReadFile("backup_test")
	require.NoError(t, err)

	backup2 := Backup{}
	backup2.FromString(string(backBuf))
	secretBuf, err := base64.URLEncoding.DecodeString(secret)
	require.NoError(t, err)
	require.False(t, bytes.Equal(backup2.SecretAsBytes(), secretBuf), "the secret may not be encrypted")

	decryptedSecret, err := decryptSecret(backup2.Secret, privateKey)
	require.NoError(t, err)
	require.True(t, bytes.Equal(secretBuf, decryptedSecret), "the secret must be decrypted")

	backupData := make(map[string]interface{})

	decryptedContentBuf, err := decrypt(backup2.Data, decryptedSecret)
	require.NoError(t, err)

	err = json.Unmarshal(decryptedContentBuf, &backupData)
	require.NoError(t, err, "decode backup failed")

	bar := input["bar"]
	bar2 := backupData["bar"]

	require.Equal(t, bar, bar2)
}

func setupTestData(t *testing.T) (map[string]interface{}, []byte) {
	input := map[string]interface{}{
		"bar": "bar string",
	}

	buf, err := json.Marshal(input)
	require.NoError(t, err)
	return input, buf
}
