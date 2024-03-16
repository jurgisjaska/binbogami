package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api/model"
)

type (
	Organization struct {
		Id          *uuid.UUID `json:"id"`
		Name        string     `json:"name"`
		Description *string    `json:"description"`
		CreatedBy   *uuid.UUID `db:"created_by" json:"createdBy"`

		CreatedAt time.Time  `db:"created_at" json:"createdAt"`
		UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt"`
	}

	Organizations []Organization

	OrganizationRepository struct {
		database *sqlx.DB
	}
)

func (r *OrganizationRepository) ById(id *uuid.UUID) (*Organization, error) {
	query := `
		SELECT o.* FROM organizations AS o
		WHERE o.id = ? AND o.deleted_at IS NULL
		LIMIT 1
	`

	organization := &Organization{}
	if err := r.database.Get(organization, query, id); err != nil {
		return nil, err
	}

	return organization, nil
}

// @todo rename this method to something better
func (r *OrganizationRepository) Find(id *uuid.UUID, member *uuid.UUID) (*Organization, error) {
	query := `
		SELECT o.* FROM organizations AS o
		JOIN members AS m ON m.organization_id = o.id
		WHERE 
		    o.id = ? AND m.user_id = ? 
		    AND m.deleted_at IS NULL AND o.deleted_at IS NULL
		LIMIT 1
	`

	organization := &Organization{}
	if err := r.database.Get(organization, query, id, member); err != nil {
		return nil, err
	}

	return organization, nil
}

func (r *OrganizationRepository) Create(o *model.Organization) (*Organization, error) {
	id := uuid.New()
	organization := &Organization{
		Id:          &id,
		Name:        o.Name,
		Description: o.Description,
		CreatedBy:   o.CreatedBy,
		CreatedAt:   time.Now(),
	}

	_, err := r.database.NamedExec(`
		INSERT INTO organizations (id, name, description, created_by, created_at)
		VALUES (:id, :name, :description, :created_by, :created_at)
	`, organization)

	if err != nil {
		return nil, err
	}

	return organization, nil
}

func (r *OrganizationRepository) ByMember(member *uuid.UUID) (*Organizations, error) {
	query := `
		SELECT o.* FROM organizations AS o
		JOIN members AS m ON m.organization_id = o.id
		WHERE 
		    m.user_id = ? 
		    AND m.deleted_at IS NULL AND o.deleted_at IS NULL
	`

	organizations := &Organizations{}
	if err := r.database.Select(organizations, query, member); err != nil {
		return nil, err
	}

	return organizations, nil
}

// CreateOrganization creates a new instance of the OrganizationRepository
func CreateOrganization(d *sqlx.DB) *OrganizationRepository {
	return &OrganizationRepository{database: d}
}
