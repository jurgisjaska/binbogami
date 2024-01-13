package model

import (
	"github.com/google/uuid"
)

type Invitation struct {
	OrganizationId *uuid.UUID `validate:"required" json:"organization_id"`
	Email          []string   `validate:"required" json:"email"`

	CreatedBy *uuid.UUID
}
