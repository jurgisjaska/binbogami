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

func (r *InvitationRepository) Create(model *model.Invitation) (Invitations, error) {
	invitations := Invitations{}
	for _, email := range model.Emails {
		id, err := uuid.NewUUID()
		if err != nil {
			return nil, err
		}

		invitation := &Invitation{
			Id:             &id,
			Email:          email,
			CreatedBy:      model.Author,
			OrganizationId: model.Organization,
			CreatedAt:      time.Now(),
			ExpiredAt:      (time.Now()).Add(defaultInvitationDuration * time.Hour),
		}

		_, err = r.database.NamedExec(`
			INSERT INTO invitations (id, email, created_by, organization_id, created_at, expired_at)
			VALUES (:id, :email, :created_by, :organization_id, :created_at, :expired_at)
		`, invitation)

		if err != nil {
			return nil, err
		}

		invitations = append(invitations, invitation)
	}

	return invitations, nil
}

func CreateInvitation(d *sqlx.DB) *InvitationRepository {
	return &InvitationRepository{database: d}
}
