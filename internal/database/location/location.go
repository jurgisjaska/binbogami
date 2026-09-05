package location

import (
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/models"
	"github.com/jurgisjaska/binbogami/internal/database/book"
)

type (
	LocationRepository interface {
		Find(id uuid.UUID) (*Location, error)
		FindMany(request *api.Request) (*Locations, int, error)
		ByBook(book *book.Book, id *uuid.UUID) (*Location, error)
		ManyByBook(book *book.Book) (*Locations, error)
		Create(l *models.Location) (*Location, error)
	}

	Location struct {
		Id          *uuid.UUID `json:"id"`
		Name        string     `json:"name"`
		Description *string    `json:"description"`
		Address     *string    `json:"address"`

		CreatedBy *uuid.UUID `db:"created_by" json:"created_by"`

		CreatedAt time.Time  `db:"created_at" json:"created_at"`
		UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
		DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
	}

	Locations []Location

	Repository struct {
		database *sqlx.DB
	}
)

// Find retrieves a Location from the repository by its ID.
func (r *Repository) Find(id uuid.UUID) (*Location, error) {
	l := &Location{}
	err := r.database.Get(l, "SELECT * FROM locations WHERE id = ? AND deleted_at IS NULL", id)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (r *Repository) FindMany(request *api.Request) (*Locations, int, error) {
	locations := &Locations{}
	q := squirrel.Select("*").From("locations").Where("deleted_at IS NULL")

	if request.Search != "" {
		q = q.Where(squirrel.Like{"name": fmt.Sprintf("%%%s%%", request.Search)})
	}

	query, args, err := q.Limit(uint64(request.Limit)).Offset(uint64(request.Offset())).ToSql()
	if err != nil {
		return nil, 0, err
	}

	err = r.database.Select(locations, query, args...)
	if err != nil {
		return nil, 0, err
	}

	query, args, err = q.RemoveColumns().Columns("COUNT(id)").ToSql()
	if err != nil {
		return nil, 0, err
	}

	var count int
	err = r.database.Get(&count, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return locations, count, nil
}

// ByBook retrieves a Location from the repository by the given CreateBook and Location IDs.
func (r *Repository) ByBook(book *book.Book, id *uuid.UUID) (*Location, error) {
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

// ManyByBook retrieves locations associated with a book.
func (r *Repository) ManyByBook(book *book.Book) (*Locations, error) {
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

func (r *Repository) Create(c *models.Location) (*Location, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	Location := &Location{
		Id:          &id,
		Name:        c.Name,
		Description: c.Description,
		Address:     c.Address,
		CreatedBy:   c.CreatedBy,
		CreatedAt:   time.Now(),
	}

	_, err = r.database.NamedExec(`
		INSERT INTO locations (id, name, description, address, created_by, created_at)
		VALUES (:id, :name, :description, :address, :created_by, :created_at)
	`, Location)

	if err != nil {
		return nil, err
	}

	return Location, nil
}

// CreateLocation creates a new instance of Repository with the specified database connection.
func CreateLocation(d *sqlx.DB) *Repository {
	return &Repository{database: d}
}
