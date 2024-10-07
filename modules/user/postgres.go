package user

import (
	"database/sql"
	"encoding/json"
	"log/slog"
)

type postgresRepository struct {
	DB *sql.DB
}

var _ Repository = (*postgresRepository)(nil)

func (p *postgresRepository) Create(params CreateUserParams) (*User, error) {
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

func (p *postgresRepository) Delete(id string) error {
	query := `
	DELETE FROM
		"gatekeeper_user"
	WHERE
		id = $1
	`

	_, err := p.DB.Exec(query, id)
	return err
}

func (p *postgresRepository) GetAll() ([]User, error) {
	query := `
	SELECT
		id,
		created_at,
		updated_at,
		deleted_at,
		login,
		password,
		extras,
		scopes
	FROM
		"gatekeeper_user"
	`

	rows, err := p.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		user, err := p.scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}

		users = append(users, *user)
	}

	return users, nil
}

type RowScanner interface {
	Scan(...any) error
}

func (postgresRepository) scanRowIntoUser(row RowScanner) (*User, error) {
	var user User
	var extrasJson, scopesJson []byte

	err := row.Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
		&user.Login,
		&user.Password,
		&extrasJson,
		&scopesJson)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(extrasJson, &user.Extras); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(scopesJson, &user.Scopes); err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *postgresRepository) GetByID(id string) (*User, error) {
	query := `
	SELECT
		id,
		created_at,
		updated_at,
		deleted_at,
		login,
		password,
		extras,
		scopes
	FROM
		"gatekeeper_user"
	WHERE
		id = $1
	`

	row := p.DB.QueryRow(query, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	return p.scanRowIntoUser(row)
}

func (p *postgresRepository) GetByLogin(login string) (*User, error) {
	query := `
	SELECT
		id,
		created_at,
		updated_at,
		deleted_at,
		login,
		password,
		extras,
		scopes
	FROM
		"gatekeeper_user"
	WHERE
		login = $1
	`

	row := p.DB.QueryRow(query, login)
	if row.Err() != nil {
		return nil, row.Err()
	}

	return p.scanRowIntoUser(row)
}

func (p *postgresRepository) Migrate(*slog.Logger) error {
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
	);
	`

	_, err := p.DB.Exec(query)
	return err
}

func (p *postgresRepository) Update(params UpdateUserParams) (*User, error) {
	query := `
	UPDATE
		"gatekeeper_user"
	SET
		updated_at = CURRENT_TIMESTAMP,
		login = $1,
		password = CASE
			WHEN $2 <> '' THEN $2
			ELSE password
		END,
		extras = $3,
		scopes = $4
	)
	WHERE
		id = $5
	`

	extrasJson, err := json.Marshal(params.Extras)
	if err != nil {
		return nil, err
	}

	scopesJson, err := json.Marshal(params.Scopes)
	if err != nil {
		return nil, err
	}

	_, err = p.DB.Exec(query, params.Login, params.Password, extrasJson, scopesJson, params.ID)
	if err != nil {
		return nil, err
	}

	return p.GetByID(params.ID)
}
