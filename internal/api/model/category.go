package model

import (
	"github.com/google/uuid"
)

type Category struct {
	Name           string     `validate:"required,gte=3" json:"name"`
	Description    *string    `json:"description"`
	OrganizationId *uuid.UUID `validate:"required" json:"organization_id"`
	CreatedBy      *uuid.UUID
}
