package book

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api"
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

	Books []Book

	Repository struct {
		database *sqlx.DB
	}
)

func (r *Repository) Create(m *model.Book) (*Book, error) {
	id := uuid.New()
	book := &Book{
		Id:             &id,
		Name:           m.Name,
		Description:    m.Description,
		OrganizationId: m.OrganizationId,
		CreatedBy:      m.CreatedBy,
		CreatedAt:      time.Now(),
	}

	_, err := r.database.NamedExec(`
		INSERT INTO books (id, name, description, organization_id, created_by, created_at)
		VALUES (:id, :name, :description, :organization_id, :created_by, :created_at)
	`, book)

	if err != nil {
		return nil, err
	}

	return book, nil
}

func (r *Repository) Find(id *uuid.UUID) (*Book, error) {
	book := &Book{}
	err := r.database.Get(book, "SELECT * FROM books WHERE id = ? AND deleted_at IS NULL", id)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (r *Repository) AddObject(book *Book, m model.BookObject) (*object, error) {
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

func (r *Repository) FindManyByOrganization(org *uuid.UUID, req *api.Request) (*Books, int, error) {
	books := &Books{}
	query := `
		SELECT * FROM books 
		WHERE organization_id = ? AND deleted_at IS NULL
		LIMIT ? OFFSET ?
	`

	offset := (req.Page - 1) * req.Limit
	err := r.database.Select(books, query, org, req.Limit, offset)
	if err != nil {
		return nil, 0, err
	}

	query = `
		SELECT COUNT(*) FROM books 
		WHERE organization_id = ? AND deleted_at IS NULL
	`
	var count int
	err = r.database.Get(&count, query, org)
	if err != nil {
		return nil, 0, err
	}

	return books, count, nil
}

func CreateBook(d *sqlx.DB) *Repository {
	return &Repository{database: d}
}
