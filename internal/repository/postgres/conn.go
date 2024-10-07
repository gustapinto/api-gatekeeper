package postgres

import (
	"database/sql"
	"errors"

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

func (Conn) InitializeDatabase(db *sql.DB) error {
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

	_, err := db.Exec(query)
	return err
}
