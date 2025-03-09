package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api/models"
	"github.com/jurgisjaska/binbogami/internal/database/book"
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

	Locations []Location

	LocationRepository struct {
		database *sqlx.DB
	}
)

// Find retrieves a Location from the repository by its ID.
func (r *LocationRepository) Find(id uuid.UUID) (*Location, error) {
	Location := &Location{}
	err := r.database.Get(Location, "SELECT * FROM locations WHERE id = ? AND deleted_at IS NULL", id.String())
	if err != nil {
		return nil, err
	}

	return Location, nil
}

// ByBook retrieves a Location from the repository by the given CreateBook and Location IDs.
func (r *LocationRepository) ByBook(book *book.Book, id *uuid.UUID) (*Location, error) {
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

// ByOrganization retrieves all locations for a given organization.
func (r *LocationRepository) ManyByOrganization(org *uuid.UUID) (*Locations, error) {
	locations := &Locations{}
	query := `
		SELECT * 
		FROM locations 
		WHERE organization_id = ?
		AND deleted_at IS NULL
	`

	err := r.database.Select(locations, query, org)
	if err != nil {
		return nil, err
	}

	return locations, nil
}

// ManyByBook retrieves locations associated with a book.
func (r *LocationRepository) ManyByBook(book *book.Book) (*Locations, error) {
	locations := &Locations{}
	query := `
		SELECT locations.* 
		FROM locations 
		JOIN books_locations AS bl ON bl.location_id = locations.id
		JOIN books AS b ON b.id = bl.book_id
		WHERE 
		    b.id = ?
		    AND locations.deleted_at IS NULL
			AND bl.deleted_at IS NULL AND b.deleted_at IS NULL
	`

	err := r.database.Select(locations, query, book.Id)
	if err != nil {
		return nil, err
	}

	return locations, nil
}

func (r *LocationRepository) Create(c *models.Location) (*Location, error) {
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
		INSERT INTO locations (id, name, description, organization_id, created_by, created_at)
		VALUES (:id, :name, :description, :organization_id, :created_by, :created_at)
	`, Location)

	if err != nil {
		return nil, err
	}

	return Location, nil
}

// CreateLocation creates a new instance of LocationRepository with the specified database connection.
func CreateLocation(d *sqlx.DB) *LocationRepository {
	return &LocationRepository{database: d}
}
