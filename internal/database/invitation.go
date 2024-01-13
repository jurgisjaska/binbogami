package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api/model"
)

const (
	defaultInvitationDuration = 24
)

type (
	// Invitation defines an entity of every invitation to join an organization send out by the email.
	// Id is used as unique key to ensure the invitation can only be used once.
	// ExpiredAt defined the invitation expiration. Every invitation should be valid for 24 hours.
	Invitation struct {
		Id             *uuid.UUID `json:"id"`
		Email          string     `json:"email"`
		CreatedBy      *uuid.UUID `db:"created_by" json:"created_by"`
		OrganizationId *uuid.UUID `db:"organization_id" json:"organization_id"`

		CreatedAt time.Time  `db:"created_at" json:"created_at"`
		OpenedAt  *time.Time `db:"opened_at" json:"opened_at"`
		DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
		ExpiredAt time.Time  `db:"expired_at" json:"expired_at"`
	}

	Invitations []*Invitation

	InvitationRepository struct {
		database *sqlx.DB
	}
)

// Open retrieves the invitation entity from the database by its UUID and marks invitation as opened.
func (r *InvitationRepository) Open(id *uuid.UUID) (*Invitation, error) {
	invitation, err := r.Find(id)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	invitation.OpenedAt = &now

	if err := r.flush(invitation); err != nil {
		return nil, err
	}

	return invitation, nil
}

// Find retrieves the invitation entity form the database by its UUID.
func (r *InvitationRepository) Find(id *uuid.UUID) (*Invitation, error) {
	invitation := &Invitation{}
	if err := r.database.Get(
		invitation,
		"SELECT * FROM invitations WHERE id = ? AND deleted_at IS NULL AND expired_at > CURRENT_TIMESTAMP()",
		id,
	); err != nil {
		return nil, err
	}

	return invitation, nil
}

func (r *InvitationRepository) Create(model *model.Invitation) (Invitations, error) {
	invitations := Invitations{}
	for _, email := range model.Email {
		id, err := uuid.NewUUID()
		if err != nil {
			return nil, err
		}

		invitation := &Invitation{
			Id:             &id,
			Email:          email,
			CreatedBy:      model.CreatedBy,
			OrganizationId: model.OrganizationId,
			CreatedAt:      time.Now(),
			ExpiredAt:      (time.Now()).Add(defaultInvitationDuration * time.Hour),
		}

		if err := r.flush(invitation); err != nil {
			return nil, err
		}

		invitations = append(invitations, invitation)
	}

	return invitations, nil
}

func (r *InvitationRepository) flush(invitation *Invitation) error {
	_, err := r.database.NamedExec(`
			INSERT INTO invitations (id, email, created_by, organization_id, created_at, opened_at, expired_at, deleted_at)
			VALUES (:id, :email, :created_by, :organization_id, :created_at, :opened_at, :expired_at, :deleted_at)
			ON DUPLICATE KEY UPDATE opened_at = :opened_at, deleted_at = :deleted_at
		`, invitation)

	if err != nil {
		return err
	}

	return nil
}

func CreateInvitation(d *sqlx.DB) *InvitationRepository {
	return &InvitationRepository{database: d}
}
