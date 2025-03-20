package service

import (
	"fmt"
	"slices"

	"github.com/gustapinto/api-gatekeeper/internal/model"
	httputil "github.com/gustapinto/api-gatekeeper/pkg/http_util"
	"golang.org/x/crypto/bcrypt"
)

type BasicAuthUserRepository interface {
	GetByLogin(string) (*model.User, error)
}

type BasicAuth struct {
	userRepository BasicAuthUserRepository
}

func NewBasicAuth(userRepository BasicAuthUserRepository) *BasicAuth {
	return &BasicAuth{
		userRepository: userRepository,
	}
}

func (s *BasicAuth) AuthenticateToken(token string) (model.User, error) {
	login, password, err := httputil.ParseBasicAuthorizationToken(token)
	if err != nil {
		return model.User{}, err
	}

	user, err := s.userRepository.GetByLogin(login)
	if err != nil {
		return model.User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return model.User{}, err
	}

	return *user, nil
}

func (s *BasicAuth) Authorize(user model.User, requiredScopes []string) error {
	for _, requiredScope := range requiredScopes {
		if !slices.Contains(user.Scopes, requiredScope) {
			return fmt.Errorf("missing %s scope", requiredScope)
		}
	}

	return nil
}
