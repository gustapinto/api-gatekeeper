package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gustapinto/api-gatekeeper/modules/gateway"
	"gopkg.in/yaml.v3"
)

type API struct {
	Address string `yaml:"address"`
}

func (a API) Validate() error {
	if strings.TrimSpace(a.Address) == "" {
		return errors.New("config 'api.address' must be present and not be empty")
	}

	return nil
}

type Database struct {
	Provider string `yaml:"provider"`
	DSN      string `yaml:"dsn"`
}

var ValidProviders = []string{
	"sqlite",
	"postgres",
}

func (d Database) Validate() error {
	if strings.TrimSpace(d.Provider) == "" {
		return errors.New("config 'database.provider' must be present and not be empty")
	}

	if !slices.Contains(ValidProviders, d.Provider) {
		return fmt.Errorf("config 'database.provider' must be one of [%s]", strings.Join(ValidProviders, ", "))
	}

	if strings.TrimSpace(d.DSN) == "" {
		return errors.New("config 'database.dsn' must be present and not be empty")
	}

	return nil
}

type Config struct {
	API      API               `yaml:"api"`
	Database Database          `yaml:"database"`
	Backends []gateway.Backend `yaml:"backends"`
}

func LoadConfig(configPath *string) (*Config, error) {
	if configPath == nil || *configPath == "" {
		return nil, errors.New("missing or empty -config=* param")
	}

	ext := strings.ToLower(filepath.Ext(*configPath))
	if ext != ".yml" && ext != ".yaml" {
		return nil, errors.New("config must have a .yml or .yaml extension")
	}

	configAbsPath, err := filepath.Abs(*configPath)
	if err != nil {
		return nil, err
	}

	configBytes, err := os.ReadFile(configAbsPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (c Config) ValidateAndNormalize() error {
	if err := c.API.Validate(); err != nil {
		return err
	}

	if err := c.Database.Validate(); err != nil {
		return err
	}

	if len(c.Backends) == 0 {
		return errors.New("config 'backends' must be present and not be empty")
	}

	for _, backend := range c.Backends {
		if err := backend.ValidateAndNormalize(); err != nil {
			return err
		}
	}

	return nil
}
