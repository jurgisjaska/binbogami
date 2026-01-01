package user

import (
	"github.com/google/uuid"
)

type (
	SetConfigurationRequest struct {
		Configuration int    `validate:"required" json:"configuration"`
		Value         string `validate:"required" json:"value"`

		UserId *uuid.UUID
	}

	ConfigurationResponse struct {
		Configuration interface{} `json:"configuration"`
	}
)
