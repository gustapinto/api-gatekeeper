package gorm

import (
	"time"

	uuidutil "github.com/gustapinto/api-gatekeeper/pkg/uuid_util"
	"gorm.io/gorm"
)

type gatekeeperUser struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Login     string `gorm:"uniqueIndex:idx_gatekeeper_user_login_uniq"`
	Password  string

	// Relationships
	Properties []gatekeeperUserProperty `gorm:"constraint:OnDelete:CASCADE"`
	Scopes     []gatekeeperUserScope    `gorm:"constraint:OnDelete:CASCADE"`
}

func (u *gatekeeperUser) BeforeSave(tx *gorm.DB) error {
	u.ID = uuidutil.NewWhenEmptyOrInvalid(u.ID)
	return nil
}

type gatekeeperUserProperty struct {
	ID               string `gorm:"primaryKey"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	GatekeeperUserID string `gorm:"uniqueIndex:idx_gatekeeper_user_properties_uniq"`
	Property         string `gorm:"uniqueIndex:idx_gatekeeper_user_properties_uniq"`
	Value            string
}

func (u *gatekeeperUserProperty) BeforeSave(tx *gorm.DB) error {
	u.ID = uuidutil.NewWhenEmptyOrInvalid(u.ID)
	return nil
}

type gatekeeperUserScope struct {
	ID               string `gorm:"primaryKey"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	GatekeeperUserID string `gorm:"uniqueIndex:idx_gatekeeper_user_scopes_uniq"`
	Scope            string `gorm:"uniqueIndex:idx_gatekeeper_user_scopes_uniq"`
}

func (u *gatekeeperUserScope) BeforeSave(tx *gorm.DB) error {
	u.ID = uuidutil.NewWhenEmptyOrInvalid(u.ID)
	return nil
}
