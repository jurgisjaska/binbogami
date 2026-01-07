package models

import (
	"github.com/google/uuid"
)

type (
	InvitationRequest struct {
		Email []string `validate:"dive,required,email" json:"email"`

		CreatedBy *uuid.UUID
	}

	InvitationResponse struct {
		Invitation interface{} `json:"invitation"`
	}
)
