package models

import (
	"github.com/google/uuid"
)

type Category struct {
	Name        string  `validate:"required,gte=3,lt=128" json:"name"`
	Description *string `json:"description"`
	Color       *string `json:"color"`

	CreatedBy *uuid.UUID
}
