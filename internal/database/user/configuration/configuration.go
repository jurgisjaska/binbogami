package configuration

import (
	"time"

	"github.com/google/uuid"
)

const defaultOrganization int = iota + 1

type (
	Configuration struct {
		Id            *uuid.UUID `json:"id"`
		Configuration int        `json:"configuration"`
		Value         string     `json:"value"`

		UserId uuid.UUID `db:"user_id" json:"userId"`

		CreatedAt time.Time  `db:"created_at" json:"createdAt"`
		UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
	}
)
