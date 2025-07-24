package crypt

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
)

// getKeyDir returns the full path to the keys directory in OS config dir
func getKeyDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "DropKey", "keys"), nil
}

func ensureKeyDirExists() error {
	keyDir, err := getKeyDir()
	if err != nil {
		return err
	}
	return os.MkdirAll(keyDir, 0o700)
}

// GenerateKey creates a new 32-byte key, stores it as base64 in a file
func GenerateKey(id string) ([]byte, error) {
	if err := ensureKeyDirExists(); err != nil {
		return nil, err
	}

	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}

	if err := SaveKey(id, key); err != nil {
		return nil, err
	}
	return key, nil
}

// SaveKey writes the provided 32-byte key (AES-256) to a base64-encoded file by ID
func SaveKey(id string, key []byte) error {
	if err := ensureKeyDirExists(); err != nil {
		return err
	}

	if len(key) != 32 {
		return errors.New("key must be 32 bytes")
	}

	keyDir, err := getKeyDir()
	if err != nil {
		return err
	}

	encoded := base64.StdEncoding.EncodeToString(key)
	keyPath := filepath.Join(keyDir, id+".key")
	return os.WriteFile(keyPath, []byte(encoded), 0o600)
}

// GetKey reads a base64-encoded key from a file by ID
func GetKey(id string) ([]byte, error) {
	keyDir, err := getKeyDir()
	if err != nil {
		return nil, err
	}

	keyPath := filepath.Join(keyDir, id+".key")
	data, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, errors.New("key not found")
	}

	key, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, errors.New("invalid key format")
	}
	if len(key) != 32 {
		return nil, errors.New("invalid key length")
	}
	return key, nil
}

// DeleteKey removes a stored key
func DeleteKey(id string) error {
	keyDir, err := getKeyDir()
	if err != nil {
		return err
	}

	keyPath := filepath.Join(keyDir, id+".key")
	return os.Remove(keyPath)
}

// MoveKey renames a key file from tempID.key to realID.key
func MoveKey(tempID, realID string) error {
	key, err := GetKey(tempID)
	if err != nil {
		return err
	}
	if err := SaveKey(realID, key); err != nil {
		return err
	}
	return DeleteKey(tempID)
}
