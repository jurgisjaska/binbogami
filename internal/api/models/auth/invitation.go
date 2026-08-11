package auth

import (
	"github.com/google/uuid"
)

type (
	InvitationRequest struct {
		Email []string `validate:"required,email" json:"email"`

		CreatedBy *uuid.UUID
	}
)
