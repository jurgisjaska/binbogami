package password

import (
	"time"

	"github.com/google/uuid"
)

// defaultPasswordResetDuration is a constant representing the password reset duration in hours
const defaultPasswordResetDuration = 2

type (
	// PasswordReset represents a data structure for storing information about a password reset.
	PasswordReset struct {
		Id     *uuid.UUID `json:"id"`
		UserId uuid.UUID  `db:"user_id" json:"userId"`

		Ip        string `db:"ip" json:"-"`
		UserAgent string `db:"user_agent" json:"-"`

		CreatedAt time.Time  `db:"created_at" json:"createdAt"`
		OpenedAt  *time.Time `db:"opened_at" json:"-"`
		ExpireAt  time.Time  `db:"expire_at" json:"-"`
	}

	PasswordResets []PasswordReset
)
