package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	yamlutil "github.com/gustapinto/api-gatekeeper/pkg/yaml_util"
)

type Config struct {
	API      API       `yaml:"api"`
	Database Database  `yaml:"database"`
	Backends []Backend `yaml:"backends"`
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

func LoadConfigFromYamlFile(configPath *string) (*Config, error) {
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
	if err := yamlutil.Unmarshal(configBytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
