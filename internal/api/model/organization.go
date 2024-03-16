package model

import (
	"github.com/google/uuid"
)

type (
	Organization struct {
		Name        string  `validate:"required,gte=3,lt=64" json:"name"`
		Description *string `validate:"required,gte=8" json:"description"`

		CreatedBy *uuid.UUID
	}
)
