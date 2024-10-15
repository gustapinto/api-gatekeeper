package postgres

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/gustapinto/api-gatekeeper/internal/model"
	_ "github.com/lib/pq"
)

type Conn struct{}

func (Conn) OpenDatabaseConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, errors.Join(errors.New("failed to connect to database"), err)
	}

	if err := db.Ping(); err != nil {
		return nil, errors.Join(errors.New("failed to connect to database"), err)
	}

	return db, nil
}

type CreateUserService interface {
	Create(model.CreateUserParams) (model.User, error)
}

func (Conn) InitializeDatabase(db *sql.DB, createUserService CreateUserService, applicationUserLogin, applicationUserPassword string) error {
	query := `
	CREATE TABLE IF NOT EXISTS "gatekeeper_user" (
		id UUID PRIMARY KEY,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP,
		deleted_at TIMESTAMP,
		login VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		extras JSONB,
		scopes JSONB
	);
	`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	_, err = createUserService.Create(model.CreateUserParams{
		Login:    applicationUserLogin,
		Password: applicationUserPassword,
		Extras:   &map[string]any{},
		Scopes: &[]string{
			"api-gatekeeper.manage-users",
		},
	})
	if err != nil {
		if strings.Contains(err.Error(), "value violates unique constraint") {
			return nil
		}

		return err
	}

	return nil
}
