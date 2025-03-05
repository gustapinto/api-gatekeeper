package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/gustapinto/api-gatekeeper/internal/config"
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

	IsAlreadyExistsError(error) bool
}

type User struct {
	userRepository UserRepository
}

func NewUser(userRepository UserRepository) User {
	return User{
		userRepository: userRepository,
	}
}

func (s User) AuthenticateToken(token string) (model.User, error) {
	if token == "" {
		return model.User{}, errors.New("badparams: missing Authorization token")
	}

	if strings.Contains(token, "Basic") {
		token = strings.TrimSpace(strings.ReplaceAll(token, "Basic", ""))
	}

	decodedToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return model.User{}, err
	}

	data := strings.Split(string(decodedToken), ":")
	if len(data) < 2 {
		return model.User{}, err
	}
	login := data[0]
	password := data[1]

	user, err := s.userRepository.GetByLogin(login)
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
		if !slices.Contains(user.Scopes, requiredScope) {
			return fmt.Errorf("missing %s scope", requiredScope)
		}
	}

	return nil
}

func (s User) Create(params model.CreateUserParams) (model.User, error) {
	if strings.TrimSpace(params.Login) == "" {
		return model.User{}, errors.New("badparams: login parameter must be present and must not be blank")
	}

	if strings.TrimSpace(params.Password) == "" {
		return model.User{}, errors.New("badparams: password parameter must be present and must not be blank")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, errors.New("badparams: failed to encode user password")
	}

	params.Password = string(hashedPassword)

	user, err := s.userRepository.Create(params)
	if err != nil {
		return model.User{}, err
	}

	user.Password = ""

	return *user, err
}

func (s User) CreateApplicationUser(cfg config.User) error {
	_, err := s.Create(model.CreateUserParams{
		Login:      cfg.Login,
		Password:   cfg.Password,
		Properties: nil,
		Scopes: []string{
			"api-gatekeeper.application",
			"api-gatekeeper.manage-users",
		},
	})
	if err != nil && !s.userRepository.IsAlreadyExistsError(err) {
		return err
	}

	return nil
}

func (s User) Update(params model.UpdateUserParams) (model.User, error) {
	if strings.TrimSpace(params.ID) == "" {
		return model.User{}, errors.New("badparams: id parameter must be present and must not be blank")
	}

	if strings.TrimSpace(params.Login) == "" {
		return model.User{}, errors.New("badparams: login parameter must be present and must not be blank")
	}

	if params.Password != nil && *params.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*params.Password), bcrypt.DefaultCost)
		if err != nil {
			return model.User{}, errors.New("badparams: failed to encode user password")
		}

		*params.Password = string(hashedPassword)
	}

	user, err := s.userRepository.Update(params)
	if err != nil {
		return model.User{}, err
	}

	user.Password = ""

	return *user, nil
}

func (s User) Delete(id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("badparams: id parameter must be present and must not be blank")
	}

	err := s.userRepository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (u User) GetByID(id string) (model.User, error) {
	if strings.TrimSpace(id) == "" {
		return model.User{}, errors.New("badparams: id parameter must be present and must not be blank")
	}

	user, err := u.userRepository.GetByID(id)
	if err != nil {
		return model.User{}, err
	}

	return *user, nil
}

func (u User) GetAll() ([]model.User, error) {
	return u.userRepository.GetAll()
}
