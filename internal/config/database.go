package config

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

type Database struct {
	Provider string `yaml:"provider"`
	DSN      string `yaml:"dsn"`
}

const (
	DatabaseProviderPostgres = "postgres"
	DatabaseProviderSqlite   = "sqlite"
)

var ValidProviders = []string{
	DatabaseProviderPostgres,
	DatabaseProviderSqlite,
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
