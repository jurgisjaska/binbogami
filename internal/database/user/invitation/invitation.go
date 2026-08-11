package invitation

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	defaultInvitationDuration = 24
)

type (
	// Invitation defines an entity of every invitation to join send out by the email.
	// Id is used as a unique key to ensure the invitation can only be used once.
	// ExpiredAt defined the invitation expiration. Every invitation should be valid for 24 hours.
	Invitation struct {
		Id        *uuid.UUID `json:"id"`
		Email     string     `json:"email"`
		Role      *int       `json:"role"`
		CreatedBy *uuid.UUID `db:"created_by" json:"created_by"`
		UserId    *uuid.UUID `db:"user_id" json:"user_id"`

		CreatedAt time.Time  `db:"created_at" json:"created_at"`
		OpenedAt  *time.Time `db:"opened_at" json:"opened_at"`
		DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
		ExpiredAt time.Time  `db:"expired_at" json:"expired_at"`
	}

	Invitations []*Invitation

	// InvitationRepository defines the interface for managing invitation entities in the database.
	InvitationRepository interface {
		Find(id uuid.UUID) (*Invitation, error)
		Update(i *Invitation) error
	}

	Repository struct {
		database *sqlx.DB
	}
)

// Find retrieves the invitation entity form the database by its UUID.
func (r *Repository) Find(id uuid.UUID) (*Invitation, error) {
	query := `
		SELECT * FROM invitations 
		WHERE id = ? AND deleted_at IS NULL AND expired_at > CURRENT_TIMESTAMP()
	`

	invitation := &Invitation{}
	if err := r.database.Get(invitation, query, id); err != nil {
		return nil, err
	}

	return invitation, nil
}

// Update updates an existing invitation record in the database with the provided invitation entity values.
func (r *Repository) Update(i *Invitation) error {
	query := `
		UPDATE invitations 
		SET 
			email = :email,
			role = :role,
			user_id = :user_id,
		    opened_at = :opened_at,
		    deleted_at = :deleted_at,
		    expired_at = :expired_at
		WHERE id = :id
	`

	_, err := r.database.NamedExec(query, i)
	if err != nil {
		return err
	}

	return nil
}

func CreateInvitation(d *sqlx.DB) *Repository {
	return &Repository{database: d}
}
