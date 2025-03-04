package gorm

import (
	"errors"
	"fmt"

	"github.com/gustapinto/api-gatekeeper/internal/config"
	"github.com/gustapinto/api-gatekeeper/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Conn struct{}

func (Conn) OpenDatabaseConnection(provider, dsn string) (db *gorm.DB, err error) {
	conf := &gorm.Config{
		FullSaveAssociations: true,
		TranslateError:       true,
	}

	switch provider {
	case config.DatabaseProviderSqlite:
		db, err = gorm.Open(sqlite.Open(dsn), conf)
		if err != nil {
			return nil, err
		}
	case config.DatabaseProviderPostgres:
		db, err = gorm.Open(postgres.Open(dsn), conf)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid provider: %s", provider)
	}

	return db, nil
}

type createUserService interface {
	Create(model.CreateUserParams) (model.User, error)
}

func (Conn) InitializeDatabase(
	db *gorm.DB,
	createUserService createUserService,
	applicationUserLogin,
	applicationUserPassword string,
) error {
	err := db.AutoMigrate(
		&gatekeeperUser{},
		&gatekeeperUserProperty{},
		&gatekeeperUserScope{})
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
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil
		}

		return err
	}

	return nil
}
