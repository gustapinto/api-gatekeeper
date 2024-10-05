package user

import (
	"database/sql"
	"encoding/json"
	"log/slog"
)

type PostgresRepository struct {
	DB *sql.DB
}

var _ Repository = (*PostgresRepository)(nil)

func (p *PostgresRepository) Create(params CreateUserParams) (*User, error) {
	query := `
	INSERT INTO "gatekeeper_user" (
		id,
		created_at,
		login,
		password,
		extras,
		scopes
	)
	VALUES (
		GEN_RANDOM_UUID(),
		CURRENT_TIMESTAMP,
		$1,
		$2,
		$3,
		$4
	)
	RETURNING id
	`

	extrasJson, err := json.Marshal(params.Extras)
	if err != nil {
		return nil, err
	}

	scopesJson, err := json.Marshal(params.Scopes)
	if err != nil {
		return nil, err
	}

	row := p.DB.QueryRow(query, params.Login, params.Password, extrasJson, scopesJson)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var id string
	if err := row.Scan(&id); err != nil {
		return nil, err
	}

	return p.GetByID(id)
}

func (p *PostgresRepository) Delete(string) error {
	panic("unimplemented")
}

func (p *PostgresRepository) GetAll() ([]User, error) {
	panic("unimplemented")
}

func (p *PostgresRepository) GetByID(string) (*User, error) {
	panic("unimplemented")
}

func (p *PostgresRepository) GetByLogin(string) (*User, error) {
	panic("unimplemented")
}

func (p *PostgresRepository) Migrate(*slog.Logger) error {
	query := `
	CREATE TABLE IF NOT EXISTS "gatekeeper_user" (
		id VARCHAR(36) UUID,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP,
		deleted_at TIMESTAMP,
		login VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		extras JSONB,
		scopes JSONB
	)
	`

	_, err := p.DB.Exec(query)
	return err
}

func (p *PostgresRepository) Update(UpdateUserParams) (*User, error) {
	panic("unimplemented")
}
