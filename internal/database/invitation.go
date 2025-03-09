package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/models"
	"github.com/jurgisjaska/binbogami/internal/database/member"
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
		CreatedBy      *uuid.UUID `db:"created_by" json:"createdBy"`
		OrganizationId *uuid.UUID `db:"organization_id" json:"organizationId"`

		CreatedAt time.Time  `db:"created_at" json:"createdAt"`
		OpenedAt  *time.Time `db:"opened_at" json:"openedAt"`
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt"`
		ExpiredAt time.Time  `db:"expired_at" json:"expiredAt"`
	}

	Invitations []*Invitation

	InvitationRepository struct {
		database *sqlx.DB
	}
)

// Open retrieves the invitation entity from the database by its UUID and marks invitation as opened.
func (r *InvitationRepository) Open(id *uuid.UUID) (*Invitation, error) {
	invitation, err := r.FindById(id)
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

// FindById retrieves the invitation entity form the database by its UUID.
func (r *InvitationRepository) FindById(id *uuid.UUID) (*Invitation, error) {
	query := `
		SELECT * FROM invitations WHERE id = ? AND deleted_at IS NULL AND expired_at > CURRENT_TIMESTAMP()
	`

	invitation := &Invitation{}
	if err := r.database.Get(invitation, query, id); err != nil {
		return nil, err
	}

	return invitation, nil
}

func (r *InvitationRepository) FindByMember(m *member.Member, req *api.Request) (*Invitations, int, error) {
	invitations := &Invitations{}

	query := `
		SELECT * 
		FROM invitations 
		WHERE organization_id = ? AND created_by = ? AND deleted_at IS NULL
		LIMIT ? OFFSET ?
	`

	offset := (req.Page - 1) * req.Limit
	err := r.database.Select(invitations, query, m.OrganizationId, m.UserId, req.Limit, offset)
	if err != nil {
		return nil, 0, err
	}

	query = `
		SELECT COUNT(*) FROM invitations 
		WHERE organization_id = ? AND created_by = ? AND deleted_at IS NULL
	`
	var count int
	err = r.database.Get(&count, query, m.OrganizationId, m.UserId)
	if err != nil {
		return nil, 0, err
	}

	return invitations, count, nil
}

func (r *InvitationRepository) Create(model *models.InvitationRequest) (Invitations, error) {
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

		if err = r.flush(invitation); err != nil {
			return nil, err
		}

		invitations = append(invitations, invitation)
	}

	return invitations, nil
}

func (r *InvitationRepository) Delete(invitation *Invitation) error {
	now := time.Now()
	invitation.DeletedAt = &now

	return r.flush(invitation)
}

func (r *InvitationRepository) flush(invitation *Invitation) error {
	query := `
		INSERT INTO invitations (id, email, created_by, organization_id, created_at, opened_at, expired_at, deleted_at)
		VALUES (:id, :email, :created_by, :organization_id, :created_at, :opened_at, :expired_at, :deleted_at)
		ON DUPLICATE KEY UPDATE opened_at = :opened_at, deleted_at = :deleted_at
	`
	_, err := r.database.NamedExec(query, invitation)
	if err != nil {
		return err
	}

	return nil
}

func CreateInvitation(d *sqlx.DB) *InvitationRepository {
	return &InvitationRepository{database: d}
}
