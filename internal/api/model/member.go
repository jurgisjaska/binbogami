package model

import (
	"github.com/google/uuid"
)

type Member struct {
	OrganizationId *uuid.UUID `validate:"required" json:"organization_id"`
	UserId         *uuid.UUID `validate:"required" json:"user_id"`
	Role           int        `validate:"required,gt=0,lt=5" json:"role"`

	CreatedBy *uuid.UUID
}
