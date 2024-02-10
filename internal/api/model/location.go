package model

import (
	"github.com/google/uuid"
)

type Location struct {
	Name        string  `validate:"required,gte=3" json:"name"`
	Description *string `json:"description"`

	OrganizationId *uuid.UUID
	CreatedBy      *uuid.UUID
}
