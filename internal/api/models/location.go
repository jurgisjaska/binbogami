package models

import (
	"github.com/google/uuid"
)

type Location struct {
	Name        string  `validate:"required,gte=3,lt=128" json:"name"`
	Description *string `json:"description"`
	Address     *string `json:"address"`

	CreatedBy *uuid.UUID
}
