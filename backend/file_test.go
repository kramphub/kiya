package backend

import (
	"bytes"
	"testing"
)

func TestEncryptDecryptSuccess(t *testing.T) {
	fileBackend := NewFileStore("./", "test") //myMasterPassword
	fileBackend.SetMasterPassword([]byte("myMasterPassword"))

	testData := []byte("testdata")
	encryptedData, err := fileBackend.encrypt(testData, fileBackend.masterPassword)
	if err != nil {
		t.Errorf("Could not encrypt data: %v", err)
	}
	if encryptedData == nil {
		t.Errorf("Expected <encrypted>: %s, got nil", testData)
	}

	decryptedData, err := fileBackend.decrypt(encryptedData, fileBackend.masterPassword)
	if err != nil {
		t.Errorf("Could not decrypt data: %v", err)
	}
	if decryptedData == nil {
		t.Errorf("Expected: %s got nil", testData)
	}
}

func TestDecryptWrongMasterPassword(t *testing.T) {
	fileBackend := NewFileStore("./", "test")
	fileBackend.SetMasterPassword([]byte("myMasterPassword"))

	testData := []byte("testdata")
	encryptedData, err := fileBackend.encrypt(testData, []byte("myMasterPassword"))
	if err != nil {
		t.Errorf("Could not encrypt data: %v", err)
	}
	if encryptedData == nil {
		t.Errorf("Expected <encrypted>: %s, got nil", testData)
	}

	_, err = fileBackend.decrypt(encryptedData, []byte("myIncorrectPassword"))
	if err == nil {
		t.Errorf("Expected: %s, got: %v", "chacha20poly1305: message authentication failed", err)
	}
}

func TestDecryptDataMismatch(t *testing.T) {
	fileBackend := NewFileStore("./", "test")
	fileBackend.SetMasterPassword([]byte("myMasterPassword"))

	testData := []byte("testdata")
	_, err := fileBackend.decrypt(testData, []byte("myIncorrectPassword"))
	if err == nil {
		t.Errorf("Expected: %s, got: %v", "data has incorrect format", err)
	}
}

func TestEncryptAlwaysDifferent(t *testing.T) {
	fileBackend := NewFileStore("./", "test")
	fileBackend.SetMasterPassword([]byte("myMasterPassword"))

	testData := []byte("testdata")
	encryptedData, _ := fileBackend.encrypt(testData, fileBackend.masterPassword)
	encryptedData2, _ := fileBackend.encrypt(testData, fileBackend.masterPassword)
	if bytes.Compare(encryptedData, encryptedData2) == 0 {
		t.Error("Expected data to be different, got equal")
	}
}

func TestNoMasterPasswordSet(t *testing.T) {
	fileBackend := NewFileStore("./", "test")

	testData := []byte("testdata")
	encryptedData, _ := fileBackend.encrypt(testData, fileBackend.masterPassword)
	encryptedData2, _ := fileBackend.encrypt(testData, fileBackend.masterPassword)
	if bytes.Compare(encryptedData, encryptedData2) == 0 {
		t.Error("Expected data to be different, got equal")
	}
}
