package entry

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/models"
)

type (
	EntryRepository interface {
		Find(id uuid.UUID) (*Entry, error)
		FindMany(request *api.Request) (*Entries, int, error)
		Create(e *models.Entry) (*Entry, error)
	}

	Entry struct {
		Id          *uuid.UUID `json:"id"`
		Amount      float64    `json:"amount"`
		Description *string    `json:"description"`

		BookId     uuid.UUID  `db:"book_id" json:"book_id"`
		CategoryId *uuid.UUID `db:"category_id" json:"category_id"`
		LocationId *uuid.UUID `db:"location_id" json:"location_id"`
		CreatedBy  *uuid.UUID `db:"created_by" json:"created_by"`

		CreatedAt time.Time  `db:"created_at" json:"created_at"`
		UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
		DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
	}

	Entries []Entry

	Repository struct {
		database *sqlx.DB
	}
)

func (r *Repository) Find(id uuid.UUID) (*Entry, error) {
	e := &Entry{}
	err := r.database.Get(e, "SELECT * FROM entries WHERE id = ? AND deleted_at IS NULL", id)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (r *Repository) FindMany(request *api.Request) (*Entries, int, error) {
	entries := &Entries{}
	query := "SELECT * FROM entries WHERE deleted_at IS NULL LIMIT ? OFFSET ?"

	err := r.database.Select(entries, query, request.Limit, request.Offset())
	if err != nil {
		return nil, 0, err
	}

	query = `SELECT COUNT(id) FROM categories WHERE deleted_at IS NULL`
	var count int
	err = r.database.Get(&count, query)
	if err != nil {
		return nil, 0, err
	}

	return entries, count, nil
}

// @deprecated
func (r *Repository) Create(e *models.Entry) (*Entry, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	entry := &Entry{
		Id:          &id,
		Amount:      e.Amount,
		Description: e.Description,
		BookId:      e.BookId,
		CategoryId:  e.CategoryId,
		LocationId:  e.LocationId,
		CreatedBy:   e.CreatedBy,
		CreatedAt:   time.Now(),
	}

	_, err = r.database.NamedExec(`
		INSERT INTO entries (id, amount, description, book_id, category_id, location_id, created_by, created_at)
		VALUES (:id, :amount, :description, :book_id, :category_id, :location_id, :created_by, :created_at)
	`, entry)

	if err != nil {
		return nil, err
	}

	return entry, nil
}

func CreateEntry(d *sqlx.DB) *Repository {
	return &Repository{database: d}
}
