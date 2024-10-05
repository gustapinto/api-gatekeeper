package gateway

import (
	"errors"
	"fmt"
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
}

func (r Route) Name() string {
	routeName := strings.ToLower(fmt.Sprintf("%s-%s", r.Method, strings.ReplaceAll(r.GatekeeperPath, "/", "-")))
	routeName = strings.ReplaceAll(routeName, "--", "-")

	return routeName
}

func (r Route) QualifiedName(qualifier string) string {
	return fmt.Sprintf("%s.%s", strings.ToLower(qualifier), r.Name())
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

type Backend struct {
	Name    string            `yaml:"name"`
	Host    string            `yaml:"host"`
	Headers map[string]string `yaml:"headers"`
	Scopes  []string          `yaml:"scopes"`
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
