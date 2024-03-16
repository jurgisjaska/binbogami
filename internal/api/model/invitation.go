package model

import (
	"github.com/google/uuid"
)

type (
	InvitationRequest struct {
		OrganizationId *uuid.UUID `validate:"required" json:"organization_id"`
		Email          []string   `validate:"required" json:"email"`

		CreatedBy *uuid.UUID
	}

	InvitationResponse struct {
		Invitation   interface{} `json:"invitation"`
		Organization interface{} `json:"organization"`
	}
)
