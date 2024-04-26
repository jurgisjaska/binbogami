package user

import (
	"github.com/google/uuid"
)

type (
	SetConfiguration struct {
		Configuration int    `validate:"required" json:"configuration"`
		Value         string `validate:"required" json:"value"`

		CreatedBy *uuid.UUID
	}
)
