package gorm

import (
	"github.com/gustapinto/api-gatekeeper/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func OpenDatabaseConnection(database config.Database) (db *gorm.DB, err error) {
	if err := database.Validate(); err != nil {
		return nil, err
	}

	var dialector gorm.Dialector
	switch database.Provider {
	case config.DatabaseProviderSqlite:
		dialector = sqlite.Open(database.DSN)
	case config.DatabaseProviderPostgres:
		dialector = postgres.Open(database.DSN)
	}

	return gorm.Open(dialector, &gorm.Config{
		FullSaveAssociations: true,
		TranslateError:       true,
		Logger:               logger.Default.LogMode(logger.Silent),
	})
}

func InitializeDatabase(
	db *gorm.DB,
) error {
	return db.AutoMigrate(
		&gatekeeperUser{},
		&gatekeeperUserProperty{},
		&gatekeeperUserScope{},
	)
}
