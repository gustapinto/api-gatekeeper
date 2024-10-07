package user

import (
	"database/sql"
	"fmt"
	"log/slog"
)

type Repository interface {
	Migrate(*slog.Logger) error

	GetAll() ([]User, error)

	GetByID(string) (*User, error)

	GetByLogin(string) (*User, error)

	Create(CreateUserParams) (*User, error)

	Update(UpdateUserParams) (*User, error)

	Delete(string) error
}

func GetRepository(config Database) (Repository, error) {
	switch config.Provider {
	case "postgres":
		db, err := sql.Open("postgres", config.DSN)
		if err != nil {
			return nil, err
		}

		if err := db.Ping(); err != nil {
			return nil, err
		}

		return &postgresRepository{
			DB: db,
		}, nil
	default:
		return nil, fmt.Errorf("provider %s is not supported", config.Provider)
	}
}
