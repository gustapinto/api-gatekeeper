package gorm

import (
	"errors"

	"github.com/gustapinto/api-gatekeeper/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	preloadClause = clause.Associations
)

type User struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) *User {
	return &User{
		db: db,
	}
}

// Public methods
func (u *User) Create(params model.CreateUserParams) (*model.User, error) {
	gUser := u.makeGatekeeperUserFromCreateUserParams(params)
	result := u.db.Create(gUser)
	if result.Error != nil {
		return nil, result.Error
	}

	return u.GetByID(gUser.ID)
}

func (u *User) Delete(userID string) error {
	result := u.db.Delete(&gatekeeperUser{}, "id = ?", userID)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (u *User) GetAll() ([]model.User, error) {
	var gUsers []gatekeeperUser
	result := u.db.Preload(preloadClause).Find(&gUsers).Order("created_at ASC")
	if result.Error != nil {
		return nil, result.Error
	}

	var users []model.User
	for _, gUser := range gUsers {
		users = append(users, *u.makeUserFromGatekeeperUser(gUser))
	}

	return users, nil
}

func (u *User) GetByID(userID string) (*model.User, error) {
	var gUser gatekeeperUser
	result := u.db.Preload(preloadClause).First(&gUser, "id = ?", userID)
	if result.Error != nil {
		return nil, result.Error
	}

	return u.makeUserFromGatekeeperUser(gUser), nil
}

func (u *User) GetByLogin(userLogin string) (*model.User, error) {
	var gUser gatekeeperUser
	result := u.db.Preload(preloadClause).First(&gUser, "login = ?", userLogin)
	if result.Error != nil {
		return nil, result.Error
	}

	return u.makeUserFromGatekeeperUser(gUser), nil
}

func (u *User) Update(params model.UpdateUserParams) (*model.User, error) {
	gUser, err := u.makeGatekeeperUserFromUpdateUserParams(params)
	if err != nil {
		return nil, err
	}

	err = u.db.Transaction(func(tx *gorm.DB) error {
		if result := tx.Delete(&gatekeeperUserProperty{}, "gatekeeper_user_id = ?", gUser.ID); result.Error != nil {
			return result.Error
		}

		if result := tx.Delete(&gatekeeperUserScope{}, "gatekeeper_user_id = ?", gUser.ID); result.Error != nil {
			return result.Error
		}

		if result := tx.Save(gUser); result.Error != nil {
			return result.Error
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return u.GetByID(gUser.ID)
}

func (*User) IsAlreadyExistsError(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, gorm.ErrDuplicatedKey)
}

// Private methods
func (*User) makeGatekeeperUserFromCreateUserParams(params model.CreateUserParams) *gatekeeperUser {
	var properties []gatekeeperUserProperty
	for property, value := range params.Properties {
		properties = append(properties, gatekeeperUserProperty{
			Property: property,
			Value:    value,
		})
	}

	var scopes []gatekeeperUserScope
	for _, scope := range params.Scopes {
		scopes = append(scopes, gatekeeperUserScope{
			Scope: scope,
		})
	}

	return &gatekeeperUser{
		Login:      params.Login,
		Password:   params.Password,
		Properties: properties,
		Scopes:     scopes,
	}
}

func (u *User) fillPropertiesForUpdate(gUser *gatekeeperUser, paramsProperties map[string]string) error {
	if gUser == nil || len(paramsProperties) == 0 {
		return nil
	}

	var dbProperties []gatekeeperUserProperty
	if err := u.db.Model(gUser).Association("Properties").Find(&dbProperties); err != nil {
		return err
	}

	for property, value := range paramsProperties {
		id := ""
		for _, dbProperty := range dbProperties {
			// Add ID to entries that already exists on the database to avoid updating the association ID
			if dbProperty.Property == property {
				id = dbProperty.ID
				break
			}
		}

		gUser.Properties = append(gUser.Properties, gatekeeperUserProperty{
			ID:               id,
			GatekeeperUserID: gUser.ID,
			Property:         property,
			Value:            value,
		})
	}

	return nil
}

func (u *User) fillScopesForUpdate(gUser *gatekeeperUser, paramsScopes []string) error {
	if gUser == nil || len(paramsScopes) == 0 {
		return nil
	}

	var dbScopes []gatekeeperUserScope
	if err := u.db.Model(gUser).Association("Scopes").Find(&dbScopes); err != nil {
		return err
	}

	for _, scope := range paramsScopes {
		id := ""
		for _, dbScope := range dbScopes {
			if dbScope.Scope == scope {
				// Add ID to entries that already exists on the database to avoid updating the association ID
				id = dbScope.ID
				break
			}
		}

		gUser.Scopes = append(gUser.Scopes, gatekeeperUserScope{
			ID:               id,
			GatekeeperUserID: gUser.ID,
			Scope:            scope,
		})
	}

	return nil
}

func (u *User) makeGatekeeperUserFromUpdateUserParams(params model.UpdateUserParams) (*gatekeeperUser, error) {
	gUser := &gatekeeperUser{
		ID:         params.ID,
		Login:      params.Login,
		Password:   *params.Password,
		Properties: []gatekeeperUserProperty{},
		Scopes:     []gatekeeperUserScope{},
	}

	if err := u.fillPropertiesForUpdate(gUser, params.Properties); err != nil {
		return nil, err
	}

	if err := u.fillScopesForUpdate(gUser, params.Scopes); err != nil {
		return nil, err
	}

	return gUser, nil
}

func (*User) makeUserFromGatekeeperUser(gUser gatekeeperUser) *model.User {
	properties := map[string]string{}
	for _, userProperty := range gUser.Properties {
		properties[userProperty.Property] = userProperty.Value
	}

	var scopes []string
	for _, userScope := range gUser.Scopes {
		scopes = append(scopes, userScope.Scope)
	}

	return &model.User{
		ID:         gUser.ID,
		Login:      gUser.Login,
		Password:   gUser.Password,
		Properties: properties,
		Scopes:     scopes,
		CreatedAt:  gUser.CreatedAt,
		UpdatedAt:  &gUser.UpdatedAt,
	}
}
