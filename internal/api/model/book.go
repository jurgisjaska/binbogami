package model

import (
	"github.com/google/uuid"
)

type Book struct {
	OrganizationId *uuid.UUID `validate:"required" json:"organization_id"`
	CreatedBy      *uuid.UUID
	Name           string  `validate:"required,gte=3" json:"name"`
	Description    *string `json:"description"`
}
