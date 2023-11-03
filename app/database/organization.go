package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type (
	Organization struct {
		Id          *uuid.UUID `json:"id"`
		Name        *string    `json:"name"`
		Description *string    `json:"description"`
		CreatedBy   *uuid.UUID `db:"created_by" json:"created_by"`
		CreatedAt   time.Time  `db:"created_at" json:"created_at"`
		UpdatedAt   *time.Time `db:"updated_at" json:"updated_at"`
		DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at"`
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
		"SELECT * FROM organizations WHERE id = ? AND deleted_at IS NULL",
		id.String(),
	); err != nil {
		return nil, err
	}

	return organization, nil
}

func CreateOrganization(d *sqlx.DB) *OrganizationRepository {
	return &OrganizationRepository{database: d}
}
