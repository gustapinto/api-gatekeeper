package gorm

import "github.com/google/uuid"

func idOrNew(id string) string {
	if len(id) == 0 {
		return uuid.NewString()
	}

	return id
}
