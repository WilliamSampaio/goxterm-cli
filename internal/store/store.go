package store

import (
	"encoding/json"
	"fmt"
	"goxterm-cli/internal/constants"
	"os"
	"path/filepath"
)

type Session struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type SshSession struct {
	Session
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type Store struct {
	Version     string       `json:"version"`
	SshSessions []SshSession `json:"ssh_sessions"`
}

func (s *Store) GetSshSession(id int) (*SshSession, error) {
	for i := range s.SshSessions {
		if s.SshSessions[i].Id == id {
			return &s.SshSessions[i], nil
		}
	}
	return nil, fmt.Errorf("SSH session with id %d not found", id)
}

func (s *Store) GetSshSessionByName(name string) (*SshSession, error) {
	for i := range s.SshSessions {
		if s.SshSessions[i].Name == name {
			return &s.SshSessions[i], nil
		}
	}
	return nil, fmt.Errorf("SSH session with name %s not found", name)
}

func (s *Store) AddSshSession(newSshSession SshSession) (*SshSession, error) {
	newId := len(s.SshSessions) + 1
	newSshSession.Id = newId
	s.SshSessions = append(s.SshSessions, newSshSession)
	return s.GetSshSession(newId)
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
		Version:     constants.AppVersion,
		SshSessions: []SshSession{},
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
