package config

import (
	"errors"
	"os"
	"path/filepath"
)

type State struct {
	Config     *Config
	ConfigPath string
}

func New() *State {
	config := NewConfig()
	state := &State{
		Config:     &config,
		ConfigPath: DefaultConfigPath,
	}

	if p := os.Getenv("PPCLI_CONFIG"); p != "" {
		state.ConfigPath = p
	}
	return state
}

var ErrorConfigPathUnknown = errors.New("config file path unknown")

func (c *State) ReadConfig() error {
	config, err := readFromFile(c.ConfigPath)
	if err != nil {
		return err
	}

	c.Config = &config
	return nil
}

func (c *State) WriteConfig() error {
	if c.ConfigPath == "" {
		return ErrorConfigPathUnknown
	}

	if c.Config == nil {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(c.ConfigPath), 0755); err != nil {
		return err
	}

	if err := c.Config.writeToFile(c.ConfigPath); err != nil {
		return err
	}
	return nil
}
