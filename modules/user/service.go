package user

import (
	"encoding/base64"
	"errors"
	"fmt"
	"slices"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Repository Repository
}

func (s Service) AuthenticateToken(token string) (User, error) {
	if token == "" {
		return User{}, errors.New("missing Authorization token")
	}

	decodedToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return User{}, err
	}

	data := strings.SplitAfter(string(decodedToken), ":")
	if len(data) < 2 {
		return User{}, err
	}
	login := data[0]
	password := data[1]

	user, err := s.Repository.GetByLogin(login)
	if err != nil {
		return User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return User{}, err
	}

	return *user, nil
}

func (s Service) Authorize(user User, requiredScopes []string) error {
	for _, requiredScope := range requiredScopes {
		if !slices.Contains(*user.Scopes, requiredScope) {
			return fmt.Errorf("missing %s scope", requiredScope)
		}
	}

	return nil
}

func (s Service) Create(params CreateUserParams) (User, error) {
	if strings.TrimSpace(params.Login) == "" {
		return User{}, errors.New("login parameter must be present and must not be blank")
	}

	if strings.TrimSpace(params.Password) == "" {
		return User{}, errors.New("password parameter must be present and must not be blank")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, errors.New("failed to encode user password")
	}

	params.Password = string(hashedPassword)

	user, err := s.Repository.Create(params)
	if err != nil {
		return User{}, err
	}

	user.Password = ""

	return *user, err
}

func (s Service) Update(params UpdateUserParams) (User, error) {
	if strings.TrimSpace(params.ID) == "" {
		return User{}, errors.New("id parameter must be present and must not be blank")
	}

	if strings.TrimSpace(params.Login) == "" {
		return User{}, errors.New("login parameter must be present and must not be blank")
	}

	if params.Password != nil && *params.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*params.Password), bcrypt.DefaultCost)
		if err != nil {
			return User{}, errors.New("failed to encode user password")
		}

		*params.Password = string(hashedPassword)
	}

	user, err := s.Repository.Update(params)
	if err != nil {
		return User{}, err
	}

	user.Password = ""

	return *user, nil
}

func (s Service) Delete(id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("id parameter must be present and must not be blank")
	}

	err := s.Repository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
