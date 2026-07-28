package auth

import (
	"github.com/google/uuid"
)

type ResetPasswordRequest struct {
	Password         string    `validate:"required,gte=8" json:"password"`
	RepeatedPassword string    `validate:"required,gte=8,eqfield=Password" json:"repeated_password"`
	Token            uuid.UUID `validate:"required,uuid" json:"token"`
}
