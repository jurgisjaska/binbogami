package database

import (
	"fmt"
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

	BookObject interface {
		Table() string
		Field() string
	}

	BookCategory struct {
		Id         int        `json:"id"`
		BookId     *uuid.UUID `db:"book_id" json:"book_id"`
		CategoryId *uuid.UUID `db:"category_id" json:"category_id"`

		CreatedBy *uuid.UUID `db:"created_by" json:"created_by"`

		CreatedAt time.Time  `db:"created_at" json:"created_at"`
		DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
	}

	BookLocation struct {
		Id         int        `json:"id"`
		BookId     *uuid.UUID `db:"book_id" json:"book_id"`
		LocationId *uuid.UUID `db:"location_id" json:"location_id"`

		CreatedBy *uuid.UUID `db:"created_by" json:"created_by"`

		CreatedAt time.Time  `db:"created_at" json:"created_at"`
		DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
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

func (r *BookRepository) Find(id *uuid.UUID) (*Book, error) {
	book := &Book{}
	err := r.database.Get(book, "SELECT * FROM books WHERE id = ? AND deleted_at IS NULL", id)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (r *BookRepository) AddObject(book *Book, m model.BookObject) (*BookObject, error) {
	e := buildObject(book, m)
	query := fmt.Sprintf(`
		INSERT INTO %s (id, book_id, %s, created_by, created_at)
		VALUES (NULL, :book_id, :%s, :created_by, :created_at)
	`, e.Table(), e.Field(), e.Field())

	_, err := r.database.NamedExec(query, e)

	if err != nil {
		return nil, err
	}

	return &e, nil
}

func buildObject(book *Book, m model.BookObject) BookObject {
	if obj, ok := m.(*model.BookCategory); ok {
		return BookCategory{
			BookId:     book.Id,
			CategoryId: obj.CategoryId,
			CreatedBy:  obj.CreatedBy,
			CreatedAt:  time.Now(),
		}
	} else if obj, ok := m.(*model.BookLocation); ok {
		return BookLocation{
			BookId:     book.Id,
			LocationId: obj.LocationId,
			CreatedBy:  obj.CreatedBy,
			CreatedAt:  time.Now(),
		}
	}

	return nil
}

func (b BookCategory) Table() string {
	return "books_categories"
}

func (b BookLocation) Table() string {
	return "books_locations"
}

func (b BookCategory) Field() string {
	return "category_id"
}

func (b BookLocation) Field() string {
	return "location_id"
}

func CreateBook(d *sqlx.DB) *BookRepository {
	return &BookRepository{database: d}
}
