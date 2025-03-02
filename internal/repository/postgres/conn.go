package postgres

import (
	"context"
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
		login VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL
	);

	CREATE TABLE IF NOT EXISTS "gatekeeper_user_property" (
		id UUID PRIMARY KEY,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP,
		gatekeeper_user_id UUID NOT NULL REFERENCES "gatekeeper_user" ("id"),
		property VARCHAR(255) NOT NULL,
		value VARCHAR(255) NOT NULL,

		UNIQUE(gatekeeper_user_id, property)
	);

	CREATE TABLE IF NOT EXISTS "gatekeeper_user_scope" (
		id UUID PRIMARY KEY,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP,
		gatekeeper_user_id UUID NOT NULL REFERENCES "gatekeeper_user" ("id"),
		scope VARCHAR(255) NOT NULL,

		UNIQUE(gatekeeper_user_id, scope)
	);
	`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	_, err = createUserService.Create(model.CreateUserParams{
		Login:      applicationUserLogin,
		Password:   applicationUserPassword,
		Properties: nil,
		Scopes: []string{
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

func Transaction(db *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
