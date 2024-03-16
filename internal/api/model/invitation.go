package model

import (
	"github.com/google/uuid"
)

type (
	InvitationRequest struct {
		// @todo organization should come from header
		OrganizationId *uuid.UUID `validate:"required" json:"organizationId"`
		Email          []string   `validate:"required" json:"email"`

		CreatedBy *uuid.UUID
	}

	InvitationResponse struct {
		Invitation   interface{} `json:"invitation"`
		Organization interface{} `json:"organization"`
	}
)
