package config

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"
)

const configDirName = "what_did_i_work_on"

type Config struct {
	Directories []DirectoryEntry `json:"directories"`
}

type DirectoryEntry struct {
	Path           string `json:"path"`
	MaxSearchDepth int    `json:"max_search_depth"`
}

func cleanPath(p string) string {
	if !path.IsAbs(p) {
		p = path.Join(p)
	}

	return filepath.Clean(p)
}

func (cfg *Config) AddPath(p string, maxSearchDepth int) {

	p = cleanPath(p)

	// check to see it already exists
	for _, d := range cfg.Directories {
		if d.Path == p {
			d.MaxSearchDepth = maxSearchDepth
			return
		}
	}

	cfg.Directories = append(cfg.Directories, DirectoryEntry{Path: p, MaxSearchDepth: maxSearchDepth})
}

func LoadConfig() (*Config, error) {

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	configDir := path.Join(userConfigDir, configDirName)

	_, err = os.Stat(configDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(configDir, 0755)
		if err != nil {
			return nil, err
		}
	}

	configFile, err := os.Open(path.Join(configDir, "config.json"))
	if err != nil {
		return &Config{}, nil
	}

	decoder := json.NewDecoder(configFile)
	cfg := &Config{}
	err = decoder.Decode(cfg)
	if err != nil {
		return &Config{}, nil
	}

	return cfg, nil

}

func SaveConfig(cfg *Config) error {
	dir, err := getConfigDir()
	if err != nil {
		return err
	}

	configFile, err := os.OpenFile(path.Join(dir, "config.json"), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer configFile.Close()

	encoder := json.NewEncoder(configFile)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}

	return nil
}

func getConfigDir() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	configDir := path.Join(userConfigDir, configDirName)

	_, err = os.Stat(configDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(configDir, 0755)
		if err != nil {
			return "", err
		}
	}

	return configDir, nil
}
