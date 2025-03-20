package config

import (
	"errors"
	"strings"
	"time"
)

type AuthType string

const (
	AuthTypeBasic AuthType = "basic"
	AuthTypeJwt   AuthType = "jwt"
)

type API struct {
	Address         string   `yaml:"address"`
	TokenExpiration string   `yaml:"tokenExpiration"`
	JwtSecret       string   `yaml:"jwtSecret"`
	AuthType        AuthType `yaml:"authType"`
	User            User     `yaml:"user"`
}

func (a API) Validate() error {
	if strings.TrimSpace(a.Address) == "" {
		return errors.New("config 'api.address' must be present and not be empty")
	}

	if strings.TrimSpace(string(a.AuthType)) == "" {
		return errors.New("config 'api.authType' must be present and not be empty")
	}

	if _, err := time.ParseDuration(a.TokenExpiration); err != nil {
		return errors.New("config 'api.tokenExpiration' must be present, not be empty and follow the https://pkg.go.dev/time#ParseDuration syntax rules")
	}

	if err := a.User.Validate(); err != nil {
		return err
	}

	return nil
}

func (a API) TokenDuration() time.Duration {
	duration, err := time.ParseDuration(a.TokenExpiration)
	if err != nil {
		duration = 30 * time.Minute
	}

	return duration
}
