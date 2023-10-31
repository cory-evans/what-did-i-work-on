package config

import (
	"encoding/json"
	"os"
	"path"
)

const configDirName = "what_did_i_work_on"

type Config struct {
	Directories []DirectoryEntry `json:"directories"`
}

type DirectoryEntry struct {
	Path           string `json:"path"`
	MaxSearchDepth int    `json:"max_search_depth"`
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

	return &Config{}, nil
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
