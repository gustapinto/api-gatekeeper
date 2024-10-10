package config

import (
	"errors"
	"net/http"
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

type apiGatekeeperUserHandler interface {
	Create(http.ResponseWriter, *http.Request)

	Update(http.ResponseWriter, *http.Request)

	Delete(http.ResponseWriter, *http.Request)

	GetByID(http.ResponseWriter, *http.Request)

	GetAll(http.ResponseWriter, *http.Request)
}

func APIGatekeeperBackend(userHandler apiGatekeeperUserHandler) Backend {
	return Backend{
		Name: "api-gatekeeper",
		Host: "",
		Scopes: []string{
			"api-gatekeeper.manage-users",
		},
		Headers: nil,
		Routes: []Route{
			{
				Method:         "POST",
				GatekeeperPath: "/api-gatekeeper/v1/users",
				HandlerFunc:    userHandler.Create,
			},
			{
				Method:         "GET",
				GatekeeperPath: "/api-gatekeeper/v1/users",
				HandlerFunc:    userHandler.GetAll,
			},
			{
				Method:         "PUT",
				GatekeeperPath: "/api-gatekeeper/v1/users/{userId}",
				HandlerFunc:    userHandler.Update,
			},
			{
				Method:         "DELETE",
				GatekeeperPath: "/api-gatekeeper/v1/users/{userId}",
				HandlerFunc:    userHandler.Delete,
			},
			{
				Method:         "GET",
				GatekeeperPath: "/api-gatekeeper/v1/users/{userId}",
				HandlerFunc:    userHandler.GetByID,
			},
		},
	}
}
