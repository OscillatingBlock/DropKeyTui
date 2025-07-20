package config

import (
	"errors"
	"os"
	"path/filepath"
)

const (
	appName     = "pasteapp"
	sessionFile = "session"
)

func getSessionPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, appName, sessionFile)
	return path, nil
}

func SaveUserID(userID string) error {
	path, err := getSessionPath()
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(path), 0o700)
	if err != nil {
		return err
	}
	return os.WriteFile(path, []byte(userID), 0o600)
}

func LoadUserID() (string, error) {
	path, err := getSessionPath()
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	userID := string(data)
	if userID == "" {
		return "", errors.New("no user ID found")
	}
	return userID, nil
}

func ClearUserID() error {
	path, err := getSessionPath()
	if err != nil {
		return err
	}
	return os.Remove(path)
}
