package service

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gustapinto/api-gatekeeper/internal/model"
)

type JWTUserRepository interface {
	GetByLogin(string) (*model.User, error)
}

type userClaims struct {
	model.User
	jwt.RegisteredClaims
}

type JWT struct {
	userRepository JWTUserRepository
	jwtSecret      string
	tokenDuration  time.Duration
}

func NewJWT(userRepository JWTUserRepository, jwtSecret string, tokenDuration time.Duration) *JWT {
	return &JWT{
		userRepository: userRepository,
		jwtSecret:      jwtSecret,
		tokenDuration:  tokenDuration,
	}
}

func (s *JWT) AuthenticateToken(token string) (model.User, error) {
	if token == "" {
		return model.User{}, errors.New("badparams: missing Authorization token")
	}

	if strings.Contains(token, "Bearer") {
		token = strings.TrimSpace(strings.ReplaceAll(token, "Bearer", ""))
	}

	t, err := jwt.ParseWithClaims(token, new(userClaims), func(t *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return model.User{}, err
	}

	if !t.Valid {
		return model.User{}, errors.New("invalid JWT token")
	}

	return t.Claims.(*userClaims).User, nil
}

func (s *JWT) Authorize(user model.User, requiredScopes []string) error {
	for _, requiredScope := range requiredScopes {
		if !slices.Contains(user.Scopes, requiredScope) {
			return fmt.Errorf("missing %s scope", requiredScope)
		}
	}

	return nil
}

func (s *JWT) GenerateToken(user model.User) (string, error) {
	claims := &userClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenDuration)),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return "Bearer " + token, nil
}
