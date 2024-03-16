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

	// SigninSuccess struct represents the successful sign-in response.
	SigninSuccess struct {
		Token string `json:"token"`
	}

	Signup struct {
		Email            *string    `validate:"required,email,lt=128" json:"email"`
		Password         string     `validate:"required,gte=8" json:"password"`
		RepeatedPassword string     `validate:"required,gte=8,eqfield=Password" json:"repeatedPassword"`
		Name             *string    `validate:"required,gte=3,lt=64" json:"name"`
		Surname          *string    `validate:"required,gte=3,lt=64" json:"surname"`
		InvitationId     *uuid.UUID `json:"invitationId"`
	}

	SignupSuccess struct {
		User  interface{} `json:"user"`
		Token string      `json:"token"`
	}
)
