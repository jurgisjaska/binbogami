package invitation

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api/models"
)

const (
	defaultInvitationDuration = 24
)

type (
	// Invitation defines an entity of every invitation to join send out by the email.
	// Id is used as unique key to ensure the invitation can only be used once.
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
		Open(id uuid.UUID) (*Invitation, error)
		Find(id uuid.UUID) (*Invitation, error)
		Create(model *models.InvitationRequest) (Invitations, error)
		Update(i *Invitation) error
		Delete(invitation *Invitation) error
	}

	Repository struct {
		database *sqlx.DB
	}
)

// Open retrieves the invitation entity from the database by its UUID and marks invitation as opened.
func (r *Repository) Open(id uuid.UUID) (*Invitation, error) {
	invitation, err := r.Find(id)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	invitation.OpenedAt = &now

	if err := r.Update(invitation); err != nil {
		return nil, err
	}

	return invitation, nil
}

// Delete marks the invitation as deleted in the database.
func (r *Repository) Delete(invitation *Invitation) error {
	now := time.Now()
	invitation.DeletedAt = &now

	return r.Update(invitation)
}

// Find retrieves the invitation entity form the database by its UUID.
func (r *Repository) Find(id uuid.UUID) (*Invitation, error) {
	query := `
		SELECT * FROM invitations WHERE id = ? AND deleted_at IS NULL AND expired_at > CURRENT_TIMESTAMP()
	`

	invitation := &Invitation{}
	if err := r.database.Get(invitation, query, id); err != nil {
		return nil, err
	}

	return invitation, nil
}

// Create generates new invitations based on the provided InvitationRequest and persists them in the database.
// @todo this should be refactored to use defined coding style
func (r *Repository) Create(model *models.InvitationRequest) (Invitations, error) {
	invitations := Invitations{}
	for _, email := range model.Email {
		id, err := uuid.NewUUID()
		if err != nil {
			return nil, err
		}

		invitation := &Invitation{
			Id:        &id,
			Email:     email,
			CreatedBy: model.CreatedBy,
			CreatedAt: time.Now(),
			ExpiredAt: (time.Now()).Add(defaultInvitationDuration * time.Hour),
		}

		if err = r.flush(invitation); err != nil {
			return nil, err
		}

		invitations = append(invitations, invitation)
	}

	return invitations, nil
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

func (r *Repository) flush(invitation *Invitation) error {
	query := `
		INSERT INTO invitations (id, email, role, created_by, user_id, created_at, opened_at, expired_at, deleted_at)
		VALUES (:id, :email, :role, :created_by, :user_id, :created_at, :opened_at, :expired_at, :deleted_at)
		ON DUPLICATE KEY UPDATE opened_at = :opened_at, deleted_at = :deleted_at, user_id = :user_id
	`
	_, err := r.database.NamedExec(query, invitation)
	if err != nil {
		return err
	}

	return nil
}

func CreateInvitation(d *sqlx.DB) *Repository {
	return &Repository{database: d}
}
