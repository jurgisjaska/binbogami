package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api/model/auth"
)

// defaultPasswordResetDuration is a constant representing the password reset duration in hours
const defaultPasswordResetDuration = 2

type (
	// PasswordReset represents a data structure for storing information about a password reset.
	PasswordReset struct {
		Id     *uuid.UUID `json:"id"`
		UserId *uuid.UUID `db:"user_id" json:"userId"`

		Ip        string `db:"ip" json:"-"`
		UserAgent string `db:"user_agent" json:"-"`

		CreatedAt time.Time  `db:"created_at" json:"createdAt"`
		OpenedAt  *time.Time `db:"opened_at" json:"-"`
		ExpireAt  time.Time  `db:"expire_at" json:"-"`
	}

	PasswordResets []PasswordReset

	// PasswordResetRepository represents a repository for storing user PasswordReset data.
	PasswordResetRepository struct {
		database *sqlx.DB
	}
)

func (r *PasswordResetRepository) FindById(id *uuid.UUID) (*PasswordReset, error) {
	query := `
		SELECT pr.* FROM user_password_resets AS pr
		WHERE pr.id = ? AND pr.expire_at > NOW()
		LIMIT 1
	`

	reset := &PasswordReset{}
	if err := r.database.Get(reset, query, id); err != nil {
		return nil, err
	}

	return reset, nil
}

func (r *PasswordResetRepository) Save(m *auth.ForgotRequest) (*PasswordReset, error) {
	id := uuid.New()
	reset := &PasswordReset{
		Id:        &id,
		UserId:    m.User.(*User).Id,
		Ip:        m.Ip,
		UserAgent: m.UserAgent,
		CreatedAt: time.Now(),
		ExpireAt:  (time.Now()).Add(defaultPasswordResetDuration * time.Hour),
	}

	_, err := r.database.NamedExec(`
		INSERT INTO user_password_resets (id, user_id, ip, user_agent, created_at, expire_at)
		VALUES (:id, :user_id, :ip, :user_agent, :created_at, :expire_at)
	`, reset)

	if err != nil {
		return nil, err
	}

	return reset, nil
}

// @todo use this insead of find by ID
func (r *PasswordResetRepository) UpdateOpenedAt(id *uuid.UUID) (*PasswordReset, error) {
	// update opened_at and return password reset token entity
	return nil, nil
}

// UpdateExpireAt updates all future user password resets with expiration date of this moment
// invalidate all password resets for the user
func (r *PasswordResetRepository) UpdateExpireAt(u *User) error {
	query := `
		UPDATE user_password_resets
		SET user_password_resets.expire_at = :expire_at
		WHERE user_password_resets.user_id = :user_id AND user_password_resets.expire_at > NOW()
	`

	_, err := r.database.NamedExec(
		query,
		map[string]interface{}{
			"expire_at": time.Now(),
			"user_id":   u.Id,
		})

	if err != nil {
		return err
	}

	return nil
}

func (r *PasswordResetRepository) FindManyByUser(u *User) (*PasswordResets, error) {
	resets := &PasswordResets{}
	query := `
		SELECT user_password_resets.* 
		FROM user_password_resets 
		WHERE user_password_resets.user_id = ? AND user_password_resets.expire_at > NOW()
	`

	err := r.database.Select(resets, query, u.Id)
	if err != nil {
		return nil, err
	}

	return resets, nil
}

// CreatePasswordReset creates a new instance of PasswordResetRepository with the specified SQL database connection.
func CreatePasswordReset(d *sqlx.DB) *PasswordResetRepository {
	return &PasswordResetRepository{database: d}
}
