package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configDirName = "what_did_i_work_on"

type Config struct {
	Directories []DirectoryEntry `json:"directories"`
}

type DirectoryEntry struct {
	Path           string `json:"path"`
	MaxSearchDepth int    `json:"max_search_depth"`
	Number         int    `json:"number"`
}

func cleanPath(p string) string {
	if !filepath.IsAbs(p) {
		p = filepath.Join(p)
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

func (cfg *Config) RemoveNumber(n int) {

	for i, d := range cfg.Directories {
		if d.Number == n {
			cfg.Directories = append(cfg.Directories[:i], cfg.Directories[i+1:]...)
			return
		}
	}
}

func LoadConfig() (*Config, error) {

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	configDir := filepath.Join(userConfigDir, configDirName)

	_, err = os.Stat(configDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(configDir, 0755)
		if err != nil {
			return nil, err
		}
	}

	configFile, err := os.Open(filepath.Join(configDir, "config.json"))
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

	configFile, err := os.OpenFile(filepath.Join(dir, "config.json"), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer configFile.Close()

	// re-number entries
	for i := 0; i < len(cfg.Directories); i++ {
		cfg.Directories[i].Number = i + 1
	}

	encoder := json.NewEncoder(configFile)
	encoder.SetIndent("", "\t")
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

	configDir := filepath.Join(userConfigDir, configDirName)

	_, err = os.Stat(configDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(configDir, 0755)
		if err != nil {
			return "", err
		}
	}

	return configDir, nil
}
