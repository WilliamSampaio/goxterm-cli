package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version   string `yaml:"version"`
	StoreType string `yaml:"store"`
	StorePath string `yaml:"store_path"`
}

func ConfigDir() string {
	dir, _ := os.UserHomeDir()
	return filepath.Join(dir, ".goxterm")
}

func configPath() string {
	return filepath.Join(ConfigDir(), "config.yaml")
}

func Exists() bool {
	_, err := os.Stat(configPath())
	return err == nil
}

func Load() (Config, error) {
	data, err := os.ReadFile(configPath())
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	return cfg, err
}

func Save(cfg Config) error {
	os.MkdirAll(filepath.Dir(configPath()), 0700)
	data, _ := yaml.Marshal(cfg)
	return os.WriteFile(configPath(), data, 0600)
}
