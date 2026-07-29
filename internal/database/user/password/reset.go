package password

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/database/user"
)

// DefaultPasswordResetDuration is a constant representing the password reset duration in hours
const DefaultPasswordResetDuration = 2

type (
	// Reset represents a data structure for storing information about a password reset.
	Reset struct {
		Id     uuid.UUID `json:"id"`
		UserId uuid.UUID `db:"user_id" json:"user_id"`

		Ip        string `db:"ip" json:"-"`
		UserAgent string `db:"user_agent" json:"-"`

		CreatedAt time.Time  `db:"created_at" json:"created_at"`
		OpenedAt  *time.Time `db:"opened_at" json:"opened_at"`
		ExpireAt  time.Time  `db:"expire_at" json:"-"`
	}

	Resets []Reset

	// PasswordResetRepository defines the interface for managing user password reset entities in the database.
	PasswordResetRepository interface {
		Create(pr *Reset) error
		Find(id uuid.UUID) (*Reset, error)
		FindManyByUser(u *user.User, limit int) (*Resets, error)
		Update(pr *Reset) error
	}

	// ResetRepository represents a repository for storing user Reset data.
	ResetRepository struct {
		database *sqlx.DB
	}
)

// Create creates a new password reset entity with the provided data.
func (r *ResetRepository) Create(pr *Reset) error {
	query := `
		INSERT INTO user_password_resets (id, user_id, ip, user_agent, created_at, expire_at)
		VALUES (:id, :user_id, :ip, :user_agent, :created_at, :expire_at)
	`

	_, err := r.database.NamedExec(query, pr)
	if err != nil {
		return err
	}

	return nil
}

// Update updates a password reset entity with the provided data.
func (r *ResetRepository) Update(pr *Reset) error {
	query := `
		UPDATE user_password_resets
		SET 
		    user_password_resets.user_id = :user_id, 
		    user_password_resets.ip = :ip, 
		    user_password_resets.user_agent = :user_agent,
		    user_password_resets.opened_at = :opened_at,
		    user_password_resets.expire_at = :expire_at
		WHERE user_password_resets.id = :id
	`

	_, err := r.database.NamedExec(query, pr)
	if err != nil {
		return err
	}

	return nil
}

// Find retrieves a password reset entity with a specific id.
func (r *ResetRepository) Find(id uuid.UUID) (*Reset, error) {
	query := `
		SELECT pr.* FROM user_password_resets AS pr
		WHERE pr.id = ? AND pr.expire_at > NOW()
		LIMIT 1
	`

	reset := &Reset{}
	if err := r.database.Get(reset, query, id); err != nil {
		return nil, err
	}

	return reset, nil
}

// FindManyByUser retrieves a list of password reset entities associated with a specific user.
func (r *ResetRepository) FindManyByUser(u *user.User, limit int) (*Resets, error) {
	query := `
		SELECT user_password_resets.* 
		FROM user_password_resets 
		WHERE user_password_resets.user_id = ? AND user_password_resets.expire_at > NOW()
		LIMIT ?
	`

	resets := &Resets{}
	err := r.database.Select(resets, query, u.Id.String(), limit+1)
	if err != nil {
		return nil, err
	}

	return resets, nil
}

// CreatePasswordReset creates a new instance of ResetRepository.
func CreatePasswordReset(d *sqlx.DB) *ResetRepository {
	return &ResetRepository{database: d}
}
