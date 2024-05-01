package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	defaultPasswordResetDuration = 24
)

type (
	PasswordReset struct {
		Id     *uuid.UUID `json:"id"`
		UserId *uuid.UUID `db:"user_id" json:"userId"`

		CreatedAt time.Time  `db:"created_at" json:"createdAt"`
		OpenedAt  *time.Time `db:"updated_at" json:"updatedAt"`
		ExpireAt  time.Time  `db:"expire_at" json:"expireAt"`
	}

	// PasswordResetRepository represents a repository for storing user PasswordReset data.
	PasswordResetRepository struct {
		database *sqlx.DB
	}
)

func (r *PasswordResetRepository) Create(user *uuid.UUID) (*PasswordReset, error) {
	id := uuid.New()
	reset := &PasswordReset{
		Id:        &id,
		UserId:    user,
		CreatedAt: time.Now(),
		ExpireAt:  (time.Now()).Add(defaultPasswordResetDuration * time.Hour),
	}

	_, err := r.database.NamedExec(`
		INSERT INTO user_password_resets (id, user_id, created_at, expire_at)
		VALUES (:id, :user_id, :created_at, :expire_at)
	`, reset)

	if err != nil {
		return nil, err
	}

	return reset, nil
}

// CreatePasswordReset creates a new instance of PasswordResetRepository with the specified SQL database connection.
func CreatePasswordReset(d *sqlx.DB) *PasswordResetRepository {
	return &PasswordResetRepository{database: d}
}
