package entry

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api/models"
)

type (
	Entry struct {
		Id          *uuid.UUID `json:"id"`
		Amount      float64    `json:"amount"`
		Description *string    `json:"description"`

		BookId     *uuid.UUID `db:"book_id" json:"book_id"`
		CategoryId *uuid.UUID `db:"category_id" json:"category_id"`
		LocationId *uuid.UUID `db:"location_id" json:"location_id"`
		CreatedBy  *uuid.UUID `db:"created_by" json:"created_by"`

		CreatedAt time.Time  `db:"created_at" json:"created_at"`
		UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
		DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
	}

	EntryRepository struct {
		database *sqlx.DB
	}
)

func (r *EntryRepository) Create(e *models.Entry) (*Entry, error) {
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

func CreateEntry(d *sqlx.DB) *EntryRepository {
	return &EntryRepository{database: d}
}
