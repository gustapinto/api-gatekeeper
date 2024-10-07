package model

import "time"

type User struct {
	ID        string          `json:"id,omitempty"`
	Login     string          `json:"login,omitempty"`
	Password  string          `json:"password,omitempty"`
	Extras    *map[string]any `json:"extras,omitempty"`
	Scopes    *[]string       `json:"scopes,omitempty"`
	CreatedAt time.Time       `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *time.Time      `json:"deleted_at,omitempty"`
}

type CreateUserParams struct {
	Login    string          `json:"login,omitempty"`
	Password string          `json:"password,omitempty"`
	Extras   *map[string]any `json:"extras,omitempty"`
	Scopes   *[]string       `json:"scopes,omitempty"`
}

type UpdateUserParams struct {
	ID       string          `json:"id,omitempty"`
	Login    string          `json:"login,omitempty"`
	Password *string         `json:"password,omitempty"`
	Extras   *map[string]any `json:"extras,omitempty"`
	Scopes   *[]string       `json:"scopes,omitempty"`
}
