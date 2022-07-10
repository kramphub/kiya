package backend

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20poly1305"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"time"
)

type FileStore struct {
	storeLocation  string
	projectID      string
	masterPassword []byte
}

func NewFileStore(storeLocation, projectID string) *FileStore {
	return &FileStore{
		projectID:     projectID,
		storeLocation: storeFileLocation(storeLocation, projectID),
	}
}

type FileStoreEntry struct {
	Value   []byte
	KeyInfo Key
}

// Get reads the store from file, fetches and decrypt the value for given key
func (f *FileStore) Get(_ context.Context, _ *Profile, key string) ([]byte, error) {
	storeData, err := f.getStore()
	if err != nil {
		return nil, err
	}

	for _, data := range storeData {
		if data.KeyInfo.Name == key {
			data, err := f.decrypt(data.Value, f.masterPassword)
			if err != nil {
				return nil, fmt.Errorf("message authentication failed")
			}
			return data, nil
		}
	}
	return nil, fmt.Errorf("%s not found", key)
}

// List reads the store from file, and fetch all keys
func (f *FileStore) List(_ context.Context, _ *Profile) (keys []Key, err error) {
	storeData, err := f.getStore()
	if err != nil {
		return nil, err
	}
	for _, info := range storeData {
		keys = append(keys, info.KeyInfo)
	}
	return keys, err
}

// CheckExists checks if given key exists in the (file)store
func (f *FileStore) CheckExists(_ context.Context, _ *Profile, key string) (bool, error) {
	storeData, err := f.getStore()
	if err != nil {
		return false, err
	}

	for _, each := range storeData {
		if each.KeyInfo.Name == key {
			return true, nil
		}
	}
	return false, nil
}

// Put a new Key with encrypted password in the store. Put overwrites the entire store file with the updated store
func (f *FileStore) Put(_ context.Context, _ *Profile, key, value string) error {
	if err := f.createStoreIfNotExists(); err != nil {
		return err
	}
	encryptedData, err := f.encrypt([]byte(value), f.masterPassword)
	if err != nil {
		return err
	}

	owner := ""
	currUser, err := user.Current()
	if err == nil {
		owner = currUser.Name
	}
	newStore := FileStoreEntry{
		Value: encryptedData,
		KeyInfo: Key{
			Name:      key,
			CreatedAt: time.Now(),
			Owner:     owner,
			Info:      "",
		},
	}

	var store []FileStoreEntry
	discStoreEntries, err := f.getStore()
	if err != nil {
		return err
	}
	if discStoreEntries != nil {
		store = append(store, discStoreEntries...)
	}
	store = append(store, newStore)
	data, err := json.Marshal(&store)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(f.storeLocation, data, 0600); err != nil {
		return err
	}
	return nil
}

// Delete a key from the store. Delete overwrites the entire store file with the updated store values
func (f *FileStore) Delete(_ context.Context, _ *Profile, key string) error {
	discStoreEntries, err := f.getStore()
	if err != nil {
		return err
	}
	var newDiscStore []FileStoreEntry
	for _, entry := range discStoreEntries {
		if entry.KeyInfo.Name != key {
			newDiscStore = append(newDiscStore, entry)
		}
	}

	data := []byte("")
	// prevents "nil" being written to file
	if len(newDiscStore) > 0 {
		data, err = json.Marshal(&newDiscStore)
		if err != nil {
			return err
		}
	}
	if err := ioutil.WriteFile(f.storeLocation, data, 0600); err != nil {
		return err
	}

	return nil
}

func (f *FileStore) Close() error {
	return nil
}

// SetMasterPassword is not relevant for this backend
func (f *FileStore) SetMasterPassword(password []byte) {
	f.masterPassword = password
}

// encrypt data based on the argon2 hashing algorithm and xchacha20 cipher algorithm
func (f *FileStore) encrypt(data, pass []byte) ([]byte, error) {
	salt := makeNonce(16)
	key := argon2.Key(pass, salt, 3, 32*1024, 4, 32)
	cipher, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}
	nonce := makeNonce(24)
	cipherText := cipher.Seal(nil, nonce, data, nil)
	return append(append(salt, nonce...), cipherText...), nil
}

// decrypt data based on the argon2 hashing algorithm and xchacha20 cipher algorithm
func (f *FileStore) decrypt(data, pass []byte) ([]byte, error) {
	if len(data) < 40 {
		return nil, errors.New("data has incorrect format")
	}
	salt := data[:16]
	nonce := data[16:40]
	data = data[40:]

	key := argon2.Key(pass, salt, 3, 32*1024, 4, 32)
	cipher, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}

	plaintext, err := cipher.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// getStore loads the file based store from disc
func (f *FileStore) getStore() ([]FileStoreEntry, error) {
	if err := f.createStoreIfNotExists(); err != nil {
		return nil, err
	}
	data, err := ioutil.ReadFile(f.storeLocation)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, nil
	}
	var store []FileStoreEntry
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}
	return store, nil
}

// createStoreIfNotExists creates the file store on disc if it does not exists and initializes with an empty value
func (f *FileStore) createStoreIfNotExists() error {
	if _, err := os.Stat(f.storeLocation); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = ioutil.WriteFile(f.storeLocation, []byte(""), 0600)
			if err != nil {
				return err
			}
		}
		return err
	}
	return nil
}

func (f *FileStore) SetParameter(key string, value interface{}) {
	if key == "masterPassword" {
		if val, ok := value.([]byte); ok {
			f.masterPassword = val
		}
	}
}

// makeNonce generates a secure random nonce used for encryption of the passwords
func makeNonce(len int) []byte {
	salt := make([]byte, len)
	n, err := rand.Reader.Read(salt)
	if err != nil {
		panic(err)
	}
	if n != len {
		panic("An error occurred while generating salt")
	}
	return salt
}

// secretStoreLocation calculates the path to the file based store
func storeFileLocation(location, projectID string) string {
	if len(location) == 0 {
		location = path.Join(os.Getenv("HOME"), fmt.Sprintf("%s.secrets.kiya", projectID))
	}
	return location
}
