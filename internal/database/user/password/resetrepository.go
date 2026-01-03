package password

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api/models/auth"
	"github.com/jurgisjaska/binbogami/internal/database/user"
)

// PasswordResetRepository represents a repository for storing user PasswordReset data.
type PasswordResetRepository struct {
	database *sqlx.DB
}

func (r *PasswordResetRepository) Save(m *auth.ForgotRequest) (*PasswordReset, error) {
	id := uuid.New()
	reset := &PasswordReset{
		Id:        &id,
		UserId:    m.User.(user.User).Id,
		Ip:        m.Ip,
		UserAgent: m.UserAgent,
		CreatedAt: time.Now(),
		ExpireAt:  (time.Now()).Add(defaultPasswordResetDuration * time.Hour),
	}

	query := `
		INSERT INTO user_password_resets (id, user_id, ip, user_agent, created_at, expire_at)
		VALUES (:id, :user_id, :ip, :user_agent, :created_at, :expire_at)
	`

	_, err := r.database.NamedExec(query, reset)
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
func (r *PasswordResetRepository) UpdateExpireAt(u *user.User) error {
	query := `
		UPDATE user_password_resets
		SET user_password_resets.expire_at = :expire_at
		WHERE user_password_resets.user_id = :user_id AND user_password_resets.expire_at > NOW()
	`

	_, err := r.database.NamedExec(query, map[string]interface{}{"expire_at": time.Now(), "user_id": u.Id})
	if err != nil {
		return err
	}

	return nil
}

// FindById retrieves a password reset entity with a specific id.
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

// FindManyByUser retrieves a list of password reset entities associated with a specific user.
func (r *PasswordResetRepository) FindManyByUser(u *user.User) (*PasswordResets, error) {
	query := `
		SELECT user_password_resets.* 
		FROM user_password_resets 
		WHERE user_password_resets.user_id = ? AND user_password_resets.expire_at > NOW()
	`

	resets := &PasswordResets{}
	err := r.database.Select(resets, query, u.Id)
	if err != nil {
		return nil, err
	}

	return resets, nil
}

// CreatePasswordReset creates a new instance of PasswordResetRepository.
func CreatePasswordReset(d *sqlx.DB) *PasswordResetRepository {
	return &PasswordResetRepository{database: d}
}
