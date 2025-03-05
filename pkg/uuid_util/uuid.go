package uuidutil

import "github.com/google/uuid"

func NewWhenEmptyOrInvalid(id string) string {
	if len(id) == 0 {
		return uuid.NewString()
	}

	if _, err := uuid.Parse(id); err != nil {
		return uuid.NewString()
	}

	return id
}
