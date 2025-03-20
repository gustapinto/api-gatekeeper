package config

import (
	"errors"
	"strings"
)

type User struct {
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
	Token    string `yaml:"token"`
}

func (u User) Validate() error {
	if strings.TrimSpace(u.Login) == "" {
		return errors.New("config 'user.login' must be present and not be empty")
	}

	if strings.TrimSpace(u.Password) == "" {
		return errors.New("config 'user.password' must be present and not be empty")
	}

	return nil
}
