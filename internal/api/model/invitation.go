package model

import (
	"github.com/google/uuid"
)

type Invitation struct {
	Organization *uuid.UUID `validate:"required" json:"organization_id"`
	Emails       []string   `validate:"required" json:"emails"`

	Author *uuid.UUID
}
