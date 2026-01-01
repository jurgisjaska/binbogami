package auth

import (
	"github.com/google/uuid"
)

type (
	// SignupRequest represents the data structure for the signup form.
	SignupRequest struct {
		Email            *string    `validate:"required,email,lt=128" json:"email"`
		Password         string     `validate:"required,gte=8" json:"password"`
		RepeatedPassword string     `validate:"required,gte=8,eqfield=Password" json:"repeatedPassword"`
		Name             *string    `validate:"required,gte=3,lt=64" json:"name"`
		Surname          *string    `validate:"required,gte=3,lt=64" json:"surname"`
		InvitationId     *uuid.UUID `json:"invitationId"`
	}

	// SignupResponse represents the response data structure for the signup process.
	SignupResponse struct {
		User  interface{} `json:"user"`
		Token string      `json:"token"`
	}
)
