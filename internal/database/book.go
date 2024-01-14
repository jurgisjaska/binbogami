package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api/model"
)

type (
	Book struct {
		Id          *uuid.UUID `json:"id"`
		Name        string     `json:"name"`
		Description *string    `json:"description"`

		OrganizationId *uuid.UUID `db:"organization_id" json:"organization_id"`
		CreatedBy      *uuid.UUID `db:"created_by" json:"created_by"`

		CreatedAt time.Time  `db:"created_at" json:"created_at"`
		UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
		DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
		ClosedAt  *time.Time `db:"closed_at" json:"closed_at"`
	}

	BookRepository struct {
		database *sqlx.DB
	}
)

func (r *BookRepository) Create(m *model.Book) (*Book, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	book := &Book{
		Id:             &id,
		Name:           m.Name,
		Description:    m.Description,
		OrganizationId: m.OrganizationId,
		CreatedBy:      m.CreatedBy,
		CreatedAt:      time.Now(),
	}

	_, err = r.database.NamedExec(`
		INSERT INTO books (id, name, description, organization_id, created_by, created_at)
		VALUES (:id, :name, :description, :organization_id, :created_by, :created_at)
	`, book)

	if err != nil {
		return nil, err
	}

	return book, nil
}

func CreateBook(d *sqlx.DB) *BookRepository {
	return &BookRepository{database: d}
}
