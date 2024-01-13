package model

import (
	"github.com/google/uuid"
	"github.com/jurgisjaska/binbogami/internal/database"
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
		Email            *string    `validate:"required,email" json:"email"`
		Password         string     `validate:"required,gte=8" json:"password"`
		RepeatedPassword string     `validate:"required,gte=8" json:"repeated_password"`
		Name             *string    `validate:"required,gte=3" json:"name"`
		Surname          *string    `validate:"required,gte=3" json:"surname"`
		Invitation       *uuid.UUID `json:"invitation"`
	}

	SignupSuccess struct {
		User  *database.User `json:"user"`
		Token string         `json:"token"`
	}
)
