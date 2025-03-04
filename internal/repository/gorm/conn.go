package gorm

import (
	"fmt"
	"strings"

	"github.com/gustapinto/api-gatekeeper/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	providerSqlite   = "sqlite"
	providerPostgres = "postgres"
)

type Conn struct{}

func (Conn) OpenDatabaseConnection(provider, dsn string) (db *gorm.DB, err error) {
	config := &gorm.Config{
		FullSaveAssociations: true,
	}

	switch provider {
	case providerSqlite:
		db, err = gorm.Open(sqlite.Open(dsn), config)
		if err != nil {
			return nil, err
		}
	case providerPostgres:
		db, err = gorm.Open(postgres.Open(dsn), config)
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
		if strings.Contains(err.Error(), "value violates unique constraint") {
			return nil
		}

		return err
	}

	return nil
}
