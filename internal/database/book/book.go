package book

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/models"
)

const (
	statusAny    string = "any"
	statusActive string = "active"
	statusClosed string = "closed"
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

// FindMany retrieves a list of books from the database based on the provided request and status.
func (r *Repository) FindMany(req *api.Request, status string) (*Books, int, error) {
	books := &Books{}
	query := fmt.Sprintf(
		`SELECT * FROM books WHERE deleted_at IS NULL %s LIMIT ? OFFSET ?`,
		r.statusQuery(status),
	)

	err := r.database.Select(books, query, req.Limit, req.Offset())
	if err != nil {
		return nil, 0, err
	}

	query = fmt.Sprintf(
		`SELECT COUNT(*) FROM books WHERE deleted_at IS NULL %s`,
		r.statusQuery(status),
	)
	var count int
	err = r.database.Get(&count, query)
	if err != nil {
		return nil, 0, err
	}

	return books, count, nil
}

func (r *Repository) statusQuery(s string) string {
	switch s {
	case statusClosed:
		return " AND closed_at IS NOT NULL "
	case statusAny:
		return ""
	case statusActive:
	default:
		return " AND closed_at IS NULL "
	}

	return ""
}

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
