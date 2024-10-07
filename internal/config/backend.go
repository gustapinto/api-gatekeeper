package config

import (
	"errors"
	"strings"
)

type Backend struct {
	Name    string            `yaml:"name"`
	Host    string            `yaml:"host"`
	Scopes  []string          `yaml:"scopes"`
	Headers map[string]string `yaml:"headers"`
	Routes  []Route           `yaml:"routes"`
}

func (b Backend) Validate() error {
	if strings.TrimSpace(b.Name) == "" {
		return errors.New("config 'backend.name' must be present and not be empty")
	}

	if strings.TrimSpace(b.Host) == "" {
		return errors.New("config 'backend.host' must be present and not be empty")
	}

	return nil
}

func (b *Backend) Normalize() {
	if b.Scopes == nil {
		b.Scopes = make([]string, 0)
	}

	if b.Headers == nil {
		b.Headers = make(map[string]string)
	}
}

func (b *Backend) ValidateAndNormalize() error {
	if err := b.Validate(); err != nil {
		return err
	}

	b.Normalize()

	for _, route := range b.Routes {
		if err := route.ValidateAndNormalize(); err != nil {
			return err
		}
	}

	return nil
}
