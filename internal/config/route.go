package config

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type Route struct {
	Method         string            `yaml:"method"`
	BackendPath    string            `yaml:"backendPath"`
	GatekeeperPath string            `yaml:"gatekeeperPath"`
	TimeoutSeconds int               `yaml:"timeoutSeconds"`
	IsPublic       bool              `yaml:"isPublic"`
	Scopes         []string          `yaml:"scopes"`
	Headers        map[string]string `yaml:"headers"`
	HandlerFunc    http.HandlerFunc
}

func (r Route) Name() string {
	routeName := strings.ToLower(fmt.Sprintf("%s-%s", r.Method, strings.ReplaceAll(r.GatekeeperPath, "/", "-")))
	routeName = strings.ReplaceAll(routeName, "--", "-")

	return routeName
}

func (r Route) Validate() error {
	if strings.TrimSpace(r.Method) == "" {
		return errors.New("config 'route.method' must be present and not be empty")
	}

	if strings.TrimSpace(r.BackendPath) == "" {
		return errors.New("config 'route.backendPath' must be present and not be empty")
	}

	if strings.TrimSpace(r.GatekeeperPath) == "" {
		return errors.New("config 'route.gatekeeperPath' must be present and not be empty")
	}

	if strings.HasPrefix(strings.ToLower(r.GatekeeperPath), "/api-gatekeeper/") {
		return errors.New("config 'route.gatekeeperPath' should not start with /api-gatekeeper, this is a reserved route namespace")
	}

	return nil
}

func (r *Route) Normalize() {
	r.Method = strings.ToUpper(r.Method)

	if r.Scopes == nil {
		r.Scopes = make([]string, 0)
	}

	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}
}

func (r *Route) ValidateAndNormalize() error {
	if err := r.Validate(); err != nil {
		return err
	}

	r.Normalize()

	return nil
}

func (r *Route) Pattern() string {
	return fmt.Sprintf("%s %s", strings.ToUpper(r.Method), r.GatekeeperPath)
}
