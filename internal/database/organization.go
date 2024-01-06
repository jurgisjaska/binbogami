package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	defaultInvitationDuration = 24
)

type (
	Organization struct {
		Id          *uuid.UUID `json:"id"`
		Name        *string    `validate:"required,gte=3,lt=64" json:"name"`
		Description *string    `validate:"required,gte=8" json:"description"`
		CreatedBy   *uuid.UUID `db:"created_by_user_id" json:"created_by"`
		OwnedBy     *uuid.UUID `db:"owned_by_user_id" json:"owned_by"`

		Members OrganizationMembers `json:"members"`

		CreatedAt time.Time  `db:"created_at" json:"created_at"`
		UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
		DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
	}

	// OrganizationUser defines an entity for Organization and User mapping table.
	// Id defines unique identifier and use numeric value as it should not be exposed.
	OrganizationUser struct {
		Id             int        `json:"id"`
		OrganizationId *uuid.UUID `db:"organization_id" json:"organization_id"`
		UserId         *uuid.UUID `db:"user_id" json:"user_id"`

		CreatedAt time.Time  `db:"created_at" json:"created_at"`
		DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
	}

	// OrganizationInvitation defines an entity of every invitation to join an organization send out by the email.
	// Id is used as unique key to ensure the invitation can only be used once.
	// ExpiredAt defined the invitation expiration. Every invitation should be valid for 24 hours.
	OrganizationInvitation struct {
		Id             *uuid.UUID
		Email          string
		CreatedBy      *uuid.UUID
		OrganizationId *uuid.UUID

		CreatedAt time.Time  `db:"created_at" json:"created_at"`
		UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
		DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
		ExpiredAt time.Time  `db:"expired_at" json:"expired_at"`
	}

	OrganizationMembers []*uuid.UUID
	Organizations       []Organization

	OrganizationRepository struct {
		database *sqlx.DB
	}
)

func (r *OrganizationRepository) Find(id uuid.UUID) (*Organization, error) {
	organization := &Organization{}
	if err := r.database.Get(
		organization,
		"SELECT organizations.* FROM organizations WHERE id = ? AND deleted_at IS NULL",
		id.String(),
	); err != nil {
		return nil, err
	}

	return organization, nil
}

func (r *OrganizationRepository) Create(org *Organization) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	org.Id = &id
	org.CreatedAt = time.Now()

	_, err = r.database.NamedExec(`
		INSERT INTO organizations (id, name, description, created_by_user_id, owned_by_user_id, created_at)
		VALUES (:id, :name, :description, :created_by_user_id, :owned_by_user_id, :created_at)
	`, org)

	if err != nil {
		return err
	}

	// @todo what happens if this fails?
	err = r.AddMember(org.Id, OrganizationMembers{org.CreatedBy})
	if err != nil {
		return err
	}

	return nil
}

func (r *OrganizationRepository) CreateInvitation(emails []string, author *uuid.UUID, org *uuid.UUID) ([]*OrganizationInvitation, error) {
	invitations := []*OrganizationInvitation{}
	for _, email := range emails {
		id, err := uuid.NewUUID()
		if err != nil {
			return nil, err
		}

		invitation := &OrganizationInvitation{
			Id:             &id,
			Email:          email,
			CreatedBy:      author,
			OrganizationId: org,
			CreatedAt:      time.Now(),
			ExpiredAt:      (time.Now()).Add(defaultInvitationDuration * time.Hour),
		}

		_, err = r.database.NamedExec(`
			INSERT INTO organizations_invitations (id, email, created_by, organization_id, created_at, expired_at)
			VALUES (:id, :email, :created_by, :organization_id, :created_at, :expired_at)
		`, org)

		if err != nil {
			return nil, err
		}

		invitations = append(invitations, invitation)
	}

	return invitations, nil
}

func (r *OrganizationRepository) AddMember(organization *uuid.UUID, members OrganizationMembers) error {

	// organizations_users should contain unique records only
	// if member was deleted but added once gain it should allow that
	// also there should not be duplicates

	for _, member := range members {
		ou := &OrganizationUser{
			OrganizationId: organization,
			UserId:         member,
			CreatedAt:      time.Now(),
		}

		_, err := r.database.NamedExec(`
			INSERT INTO organizations_users (id, organization_id, user_id, created_at)
			VALUES (NULL, :organization_id, :user_id, :created_at)
		`, ou)

		if err != nil {
			return err
		}
	}

	return nil
}

func CreateOrganization(d *sqlx.DB) *OrganizationRepository {
	return &OrganizationRepository{database: d}
}
