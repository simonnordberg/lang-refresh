package config

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
)

var DefaultConfigPath string

func init() {
	usr, err := user.Current()
	if err != nil {
		return
	}

	if usr.HomeDir != "" {
		DefaultConfigPath = filepath.Join(usr.HomeDir, ".config", "ppcli", "config.json")
	}
}

type Config struct {
	Contexts []*ConfigContext
}

func NewConfig() Config {
	return Config{
		Contexts: make([]*ConfigContext, 0),
	}
}

type ConfigContext struct {
	Name  string
	Token string
	Url   string
}

func (config *Config) ContextNames() []string {
	if len(config.Contexts) == 0 {
		return nil
	}
	names := make([]string, len(config.Contexts))
	for i, ctx := range config.Contexts {
		names[i] = ctx.Name
	}
	return names
}

func (config *Config) ContextByName(name string) *ConfigContext {
	for _, c := range config.Contexts {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func (config *Config) writeToFile(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(config); err != nil {
		return err
	}

	return nil
}

func readFromFile(filePath string) (Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return NewConfig(), err
	}
	defer file.Close()

	config := NewConfig()
	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		return NewConfig(), err
	}

	return config, nil
}
