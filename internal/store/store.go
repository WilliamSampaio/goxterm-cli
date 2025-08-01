package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type SshCredential struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type Store struct {
	Credentials map[string]SshCredential `json:"credentials"`
}

func Exists(storePath string) bool {
	_, err := os.Stat(storePath)
	return err == nil
}

func Initialize(storePath string) error {
	if Exists(storePath) {
		return fmt.Errorf("store already exists at %s", storePath)
	}

	store := Store{
		Credentials: make(map[string]SshCredential),
	}

	return Save(storePath, &store)
}

func Load(storePath string) (Store, error) {
	data, err := os.ReadFile(storePath)
	if err != nil {
		return Store{}, err
	}

	var store Store
	err = json.Unmarshal(data, &store)
	return store, err
}

func Save(storePath string, store *Store) error {
	os.MkdirAll(filepath.Dir(storePath), 0700)

	data, err := json.Marshal(store)
	if err != nil {
		return fmt.Errorf("failed to initialize store: %v", err)
	}

	return os.WriteFile(storePath, data, 0600)
}
