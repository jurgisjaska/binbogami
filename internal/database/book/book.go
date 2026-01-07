package book

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/models"
)

type (
	Book struct {
		Id          *uuid.UUID `json:"id"`
		Name        string     `json:"name"`
		Description *string    `json:"description"`

		CreatedBy *uuid.UUID `db:"created_by" json:"createdBy"`

		CreatedAt time.Time  `db:"created_at" json:"createdAt"`
		UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt"`
		ClosedAt  *time.Time `db:"closed_at" json:"closedAt"`
	}

	Books []Book

	Repository struct {
		database *sqlx.DB
	}
)

func (r *Repository) Create(m *models.CreateBook) (*Book, error) {
	id := uuid.New()
	book := &Book{
		Id:          &id,
		Name:        m.Name,
		Description: m.Description,
		CreatedBy:   m.CreatedBy,
		CreatedAt:   time.Now(),
	}

	_, err := r.database.NamedExec(`
		INSERT INTO books (id, name, description, created_by, created_at)
		VALUES (:id, :name, :description, :created_by, :created_at)
	`, book)

	if err != nil {
		return nil, err
	}

	return book, nil
}

func (r *Repository) Update(e *Book, m *models.UpdateBook) (*Book, error) {
	e.Name = m.Name
	e.Description = m.Description

	_, err := r.database.NamedExec(`
		UPDATE books
		SET name = :name, description = :description
		WHERE id = :id AND deleted_at IS NULL
	`, e)

	if err != nil {
		return nil, err
	}

	return e, nil
}

// Find retrieves a book by its ID from the database if it exists and hasn't been marked as deleted.
func (r *Repository) Find(id *uuid.UUID) (*Book, error) {
	book := &Book{}
	err := r.database.Get(book, "SELECT * FROM books WHERE id = ? AND deleted_at IS NULL", id)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (r *Repository) AddObject(book *Book, m models.BookObject) (*object, error) {
	e := buildObject(book, m)
	query := fmt.Sprintf(`
		INSERT INTO %s (id, book_id, %s, created_by, created_at)
		VALUES (NULL, :book_id, :%s, :created_by, :created_at)
	`, e.table(), e.field(), e.field())

	_, err := r.database.NamedExec(query, e)

	if err != nil {
		return nil, err
	}

	return &e, nil
}

func CreateBook(d *sqlx.DB) *Repository {
	return &Repository{database: d}
}
