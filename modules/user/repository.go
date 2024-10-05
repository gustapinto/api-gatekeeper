package user

import (
	"log/slog"
	"time"
)

type User struct {
	ID        string
	Login     string
	Password  string
	Extras    *map[string]any
	Scopes    *[]string
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

type CreateUserParams struct {
	Login    string
	Password string
	Extras   *map[string]any
	Scopes   *[]string
}

type UpdateUserParams struct {
	ID       string
	Login    string
	Password string
	Extras   *map[string]any
	Scopes   *[]string
}

type Repository interface {
	Migrate(*slog.Logger) error

	GetAll() ([]User, error)

	GetByID(string) (*User, error)

	GetByLogin(string) (*User, error)

	Create(CreateUserParams) (*User, error)

	Update(UpdateUserParams) (*User, error)

	Delete(string) error
}
