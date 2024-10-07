package config

import (
	"errors"
	"strings"
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
