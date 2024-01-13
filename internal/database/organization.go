package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type (
	Organization struct {
		Id          *uuid.UUID `json:"id"`
		Name        *string    `validate:"required,gte=3,lt=64" json:"name"`
		Description *string    `validate:"required,gte=8" json:"description"`
		CreatedBy   *uuid.UUID `db:"created_by" json:"created_by"`

		CreatedAt time.Time  `db:"created_at" json:"created_at"`
		UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
		DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
	}

	Organizations []Organization

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
		INSERT INTO organizations (id, name, description, created_by, created_at)
		VALUES (:id, :name, :description, :created_by, :created_at)
	`, org)

	if err != nil {
		return err
	}

	return nil
}

func CreateOrganization(d *sqlx.DB) *OrganizationRepository {
	return &OrganizationRepository{database: d}
}
