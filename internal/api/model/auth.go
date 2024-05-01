package model

import (
	"github.com/google/uuid"
)

type (
	// Signin struct is a data structure representing the user sign-in information.
	Signin struct {
		Email    string `validate:"required,email" json:"email"`
		Password string `validate:"required" json:"password"`
	}

	// SigninResponse struct represents the successful sign-in response.
	SigninResponse struct {
		Token        string      `json:"token"`
		User         interface{} `json:"user"`
		Organization interface{} `json:"organization"`
		Member       bool        `json:"member"`
	}

	SignupRequest struct {
		Email            *string    `validate:"required,email,lt=128" json:"email"`
		Password         string     `validate:"required,gte=8" json:"password"`
		RepeatedPassword string     `validate:"required,gte=8,eqfield=Password" json:"repeatedPassword"`
		Name             *string    `validate:"required,gte=3,lt=64" json:"name"`
		Surname          *string    `validate:"required,gte=3,lt=64" json:"surname"`
		InvitationId     *uuid.UUID `json:"invitationId"`
	}

	SignupResponse struct {
		User         interface{} `json:"user"`
		Token        string      `json:"token"`
		Member       bool        `json:"member"`
		Organization interface{} `json:"organization"`
	}

	// ForgotPasswordRequest struct is a data structure representing the request for forgot password feature,
	ForgotPasswordRequest struct {
		Email string `validate:"required,email" json:"email"`
	}

	ResetPasswordRequest struct {
		Password         string `validate:"required,gte=8" json:"password"`
		RepeatedPassword string `validate:"required,gte=8,eqfield=Password" json:"repeatedPassword"`
	}
)
