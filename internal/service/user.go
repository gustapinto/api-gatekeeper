package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/gustapinto/api-gatekeeper/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetAll() ([]model.User, error)

	GetByID(string) (*model.User, error)

	GetByLogin(string) (*model.User, error)

	Create(model.CreateUserParams) (*model.User, error)

	Update(model.UpdateUserParams) (*model.User, error)

	Delete(string) error
}

type User struct {
	Repository UserRepository
}

func (s User) AuthenticateToken(token string) (model.User, error) {
	if token == "" {
		return model.User{}, errors.New("missing Authorization token")
	}

	decodedToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return model.User{}, err
	}

	data := strings.SplitAfter(string(decodedToken), ":")
	if len(data) < 2 {
		return model.User{}, err
	}
	login := data[0]
	password := data[1]

	user, err := s.Repository.GetByLogin(login)
	if err != nil {
		return model.User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return model.User{}, err
	}

	return *user, nil
}

func (s User) Authorize(user model.User, requiredScopes []string) error {
	for _, requiredScope := range requiredScopes {
		if !slices.Contains(*user.Scopes, requiredScope) {
			return fmt.Errorf("missing %s scope", requiredScope)
		}
	}

	return nil
}

func (s User) Create(params model.CreateUserParams) (model.User, error) {
	if strings.TrimSpace(params.Login) == "" {
		return model.User{}, errors.New("login parameter must be present and must not be blank")
	}

	if strings.TrimSpace(params.Password) == "" {
		return model.User{}, errors.New("password parameter must be present and must not be blank")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, errors.New("failed to encode user password")
	}

	params.Password = string(hashedPassword)

	user, err := s.Repository.Create(params)
	if err != nil {
		return model.User{}, err
	}

	user.Password = ""

	return *user, err
}

func (s User) Update(params model.UpdateUserParams) (model.User, error) {
	if strings.TrimSpace(params.ID) == "" {
		return model.User{}, errors.New("id parameter must be present and must not be blank")
	}

	if strings.TrimSpace(params.Login) == "" {
		return model.User{}, errors.New("login parameter must be present and must not be blank")
	}

	if params.Password != nil && *params.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*params.Password), bcrypt.DefaultCost)
		if err != nil {
			return model.User{}, errors.New("failed to encode user password")
		}

		*params.Password = string(hashedPassword)
	}

	user, err := s.Repository.Update(params)
	if err != nil {
		return model.User{}, err
	}

	user.Password = ""

	return *user, nil
}

func (s User) Delete(id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("id parameter must be present and must not be blank")
	}

	err := s.Repository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
