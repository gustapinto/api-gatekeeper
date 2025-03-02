package postgres

import (
	"context"
	"database/sql"
	"slices"

	"github.com/gustapinto/api-gatekeeper/internal/model"
)

type User struct {
	db *sql.DB
}

func NewUser(db *sql.DB) User {
	return User{
		db: db,
	}
}

func (p User) upsertUserProperty(
	tx *sql.Tx,
	userID string,
	property string,
	value string,
) error {
	query := `
	INSERT INTO "gatekeeper_user_property" (
		id,
		created_at,
		gatekeeper_user_id,
		property,
		value
	) VALUES (
		GEN_RANDOM_UUID(),
		CURRENT_TIMESTAMP,
		$1,
		$2,
		$3::VARCHAR
	) ON CONFLICT (
		gatekeeper_user_id,
		property
	) DO UPDATE SET
	 	updated_at = EXCLUDED.created_at,
		value = EXCLUDED.value
	`

	_, err := tx.ExecContext(
		context.Background(),
		query,
		userID,
		property,
		value)
	if err != nil {
		return err
	}

	return nil
}

func (p User) upsertUserScope(tx *sql.Tx, userID string, scope string) error {
	query := `
	INSERT INTO "gatekeeper_user_scope" (
		id,
		created_at,
		gatekeeper_user_id,
		scope
	) VALUES (
		GEN_RANDOM_UUID(),
		CURRENT_TIMESTAMP,
		$1,
		$2
	) ON CONFLICT (
		gatekeeper_user_id,
		scope
	) DO UPDATE SET
	 	updated_at = EXCLUDED.created_at
	`

	_, err := tx.ExecContext(
		context.Background(),
		query,
		userID,
		scope)
	if err != nil {
		return err
	}

	return nil
}

func (p User) insertUser(tx *sql.Tx, params model.CreateUserParams) (string, error) {
	query := `
	INSERT INTO "gatekeeper_user" (
		id,
		created_at,
		login,
		password
	) VALUES (
		GEN_RANDOM_UUID(),
		CURRENT_TIMESTAMP,
		$1,
		$2
	)
	RETURNING id
	`

	row := tx.QueryRow(query, params.Login, params.Password)
	if row.Err() != nil {
		return "", row.Err()
	}

	var id string
	if err := row.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}

func (p User) Create(params model.CreateUserParams) (_ *model.User, err error) {
	var id string

	err = Transaction(p.db, func(tx *sql.Tx) error {
		id, err = p.insertUser(tx, params)
		if err != nil {
			return err
		}

		for property, value := range params.Properties {
			if err := p.upsertUserProperty(tx, id, property, value); err != nil {
				return err
			}
		}

		for _, scope := range params.Scopes {
			if err := p.upsertUserScope(tx, id, scope); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return p.GetByID(id)
}

func (User) deleteUserProperties(tx *sql.Tx, userID string) error {
	query := `
	DELETE FROM
		"gatekeeper_user_property"
	WHERE
		gatekeeper_user_id = $1
	`

	_, err := tx.ExecContext(context.Background(), query, userID)
	return err
}

func (User) deleteUserProperty(tx *sql.Tx, property, userID string) error {
	query := `
	DELETE FROM
		"gatekeeper_user_property"
	WHERE
		gatekeeper_user_id = $1
		AND property = $2
	`

	_, err := tx.ExecContext(context.Background(), query, userID, property)
	return err
}

func (User) deleteUserScopes(tx *sql.Tx, userID string) error {
	query := `
	DELETE FROM
		"gatekeeper_user_scope"
	WHERE
		gatekeeper_user_id = $1
	`

	_, err := tx.ExecContext(context.Background(), query, userID)
	return err
}

func (User) deleteUserScope(tx *sql.Tx, scope, userID string) error {
	query := `
	DELETE FROM
		"gatekeeper_user_scope"
	WHERE
		gatekeeper_user_id = $1
		AND scope = $2
	`

	_, err := tx.ExecContext(context.Background(), query, userID, scope)
	return err
}

func (User) deleteUser(tx *sql.Tx, userID string) error {
	query := `
	DELETE FROM
		"gatekeeper_user"
	WHERE
		id = $1
	`

	_, err := tx.Exec(query, userID)
	return err
}

func (p User) Delete(id string) error {
	return Transaction(p.db, func(tx *sql.Tx) error {
		if err := p.deleteUserProperties(tx, id); err != nil {
			return err
		}

		if err := p.deleteUserScopes(tx, id); err != nil {
			return err
		}

		if err := p.deleteUser(tx, id); err != nil {
			return err
		}

		return nil
	})
}

func (p User) getScopesForUser(tx *sql.Tx, userID string) ([]string, error) {
	query := `
	SELECT
		scope
	FROM
		"gatekeeper_user_scope"
	WHERE
		gatekeeper_user_id = $1
	ORDER BY
		scope ASC
	`

	rows, err := tx.QueryContext(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scopes []string
	for rows.Next() {
		var scope string
		if err := rows.Scan(&scope); err != nil {
			return nil, err
		}

		scopes = append(scopes, scope)
	}

	return scopes, nil
}

func (p User) getPropertiesForUser(tx *sql.Tx, userID string) (map[string]string, error) {
	query := `
	SELECT
		property,
		value
	FROM
		"gatekeeper_user_property"
	WHERE
		gatekeeper_user_id = $1
	ORDER BY
		property ASC
	`

	rows, err := tx.QueryContext(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	properties := make(map[string]string)
	for rows.Next() {
		var (
			property string
			value    string
		)
		if err := rows.Scan(&property, &value); err != nil {
			return nil, err
		}

		properties[property] = value
	}

	return properties, nil
}

func (p User) fillUserWithPropertiesAndScopes(tx *sql.Tx, user *model.User) (err error) {
	user.Scopes, err = p.getScopesForUser(tx, user.ID)
	if err != nil {
		return err
	}

	user.Properties, err = p.getPropertiesForUser(tx, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (p User) GetAll() (users []model.User, err error) {
	query := `
	SELECT
		id,
		created_at,
		updated_at,
		login,
		'' as password
	FROM
		"gatekeeper_user"
	`

	err = Transaction(p.db, func(tx *sql.Tx) error {
		rows, err := p.db.Query(query)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			user, err := p.scanRowIntoUser(rows)
			if err != nil {
				return err
			}

			if err := p.fillUserWithPropertiesAndScopes(tx, user); err != nil {
				return err
			}

			users = append(users, *user)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return users, nil
}

type RowScanner interface {
	Scan(...any) error
}

func (User) scanRowIntoUser(row RowScanner) (*model.User, error) {
	var user model.User

	err := row.Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Login,
		&user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p User) GetByID(id string) (user *model.User, err error) {
	query := `
	SELECT
		id,
		created_at,
		updated_at,
		login,
		password
	FROM
		"gatekeeper_user"
	WHERE
		id = $1
	`

	err = Transaction(p.db, func(tx *sql.Tx) error {
		row := tx.QueryRow(query, id)
		if row.Err() != nil {
			return row.Err()
		}

		user, err = p.scanRowIntoUser(row)
		if err != nil {
			return err
		}

		if err := p.fillUserWithPropertiesAndScopes(tx, user); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (p User) GetByLogin(login string) (user *model.User, err error) {
	query := `
	SELECT
		id,
		created_at,
		updated_at,
		login,
		password
	FROM
		"gatekeeper_user"
	WHERE
		login = $1
	`

	err = Transaction(p.db, func(tx *sql.Tx) error {
		row := tx.QueryRow(query, login)
		if row.Err() != nil {
			return row.Err()
		}

		user, err = p.scanRowIntoUser(row)
		if err != nil {
			return err
		}

		if err := p.fillUserWithPropertiesAndScopes(tx, user); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (p User) Update(params model.UpdateUserParams) (*model.User, error) {
	query := `
	UPDATE
		"gatekeeper_user"
	SET
		updated_at = CURRENT_TIMESTAMP,
		login = $1,
		password = (CASE
			WHEN $2 <> '' THEN $2
			ELSE password
		END)
	WHERE
		id = $3
	`

	err := Transaction(p.db, func(tx *sql.Tx) error {
		existingProperties, err := p.getPropertiesForUser(tx, params.ID)
		if err != nil {
			return err
		}

		for property := range existingProperties {
			if _, ok := params.Properties[property]; !ok {
				if err := p.deleteUserProperty(tx, property, params.ID); err != nil {
					return err
				}
			}
		}

		existingScopes, err := p.getScopesForUser(tx, params.ID)
		if err != nil {
			return err
		}

		for _, scope := range existingScopes {
			if !slices.Contains(params.Scopes, scope) {
				if err := p.deleteUserScope(tx, scope, params.ID); err != nil {
					return err
				}
			}
		}

		_, err = tx.ExecContext(
			context.Background(),
			query,
			params.Login,
			params.Password,
			params.ID)
		if err != nil {
			return err
		}

		for property, value := range params.Properties {
			if err := p.upsertUserProperty(tx, params.ID, property, value); err != nil {
				return err
			}
		}

		for _, scope := range params.Scopes {
			if err := p.upsertUserScope(tx, params.ID, scope); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return p.GetByID(params.ID)
}
