package model

import (
	"github.com/google/uuid"
)

// @todo add name length limits (the same length as database field length)

type Category struct {
	Name        string  `validate:"required,gte=3" json:"name"`
	Description *string `json:"description"`

	OrganizationId *uuid.UUID
	CreatedBy      *uuid.UUID
}
