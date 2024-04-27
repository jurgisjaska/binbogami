package user

import (
	"github.com/google/uuid"
)

type (
	SetConfiguration struct {
		Configuration int    `validate:"required" json:"configuration"`
		Value         string `validate:"required" json:"value"`

		UserId *uuid.UUID
	}

	ConfigurationResponse struct {
		Configuration interface{} `json:"configuration"`
		Organization  interface{} `json:"organization"`
	}
)
