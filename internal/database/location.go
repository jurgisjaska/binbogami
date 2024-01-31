package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api/model"
)

type (
	Location struct {
		Id          *uuid.UUID `json:"id"`
		Name        string     `json:"name"`
		Description *string    `json:"description"`

		OrganizationId *uuid.UUID `db:"organization_id" json:"organization_id"`
		CreatedBy      *uuid.UUID `db:"created_by" json:"created_by"`

		CreatedAt time.Time  `db:"created_at" json:"created_at"`
		UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
		DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
	}

	LocationRepository struct {
		database *sqlx.DB
	}
)

func (r *LocationRepository) Find(id uuid.UUID) (*Location, error) {
	Location := &Location{}
	err := r.database.Get(Location, "SELECT * FROM locations WHERE id = ? AND deleted_at IS NULL", id.String())
	if err != nil {
		return nil, err
	}

	return Location, nil
}

func (r *LocationRepository) ByBook(book *Book, id *uuid.UUID) (*Location, error) {
	query := `
		SELECT locations.* 
		FROM locations 
		JOIN books_locations AS bl ON bl.location_id = locations.id
		JOIN books AS b ON b.id = bl.book_id
		WHERE 
		    b.id = ? AND locations.id = ?
		    AND locations.deleted_at IS NULL
			AND bl.deleted_at IS NULL AND b.deleted_at IS NULL
	`

	location := &Location{}
	if err := r.database.Get(location, query, book.Id, id); err != nil {
		return nil, err
	}

	return location, nil
}

func (r *LocationRepository) Create(c *model.Location) (*Location, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	Location := &Location{
		Id:             &id,
		Name:           c.Name,
		Description:    c.Description,
		OrganizationId: c.OrganizationId,
		CreatedBy:      c.CreatedBy,
		CreatedAt:      time.Now(),
	}

	_, err = r.database.NamedExec(`
		INSERT INTO locations (id, name, description, organization_id, created_by, created_at, updated_at, deleted_at)
		VALUES (:id, :name, :description, :organization_id, :created_by, :created_at, :updated_at, :deleted_at)
	`, Location)

	if err != nil {
		return nil, err
	}

	return Location, nil
}

func CreateLocation(d *sqlx.DB) *LocationRepository {
	return &LocationRepository{database: d}
}
