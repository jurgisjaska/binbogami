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
	BookRepository interface {
		Find(id uuid.UUID) (*Book, error)
		FindMany(request *api.Request, status string) (*Books, int, error)
		FindManyByName(req *api.Request, status string, search string) (*Books, int, error)
		Create(book *Book) error
		Update(book *Book) error
		AddObject(book *Book, model models.BookObject) (any, error)
	}

	Book struct {
		Id          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Description *string   `json:"description"`

		CreatedBy uuid.UUID `db:"created_by" json:"created_by"`
		Author    *string   `db:"-" json:"author"`

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

func (r *Repository) Find(id uuid.UUID) (*Book, error) {
	book := &Book{}
	err := r.database.Get(book, "SELECT * FROM books WHERE id = ? AND deleted_at IS NULL", id)
	if err != nil {
		return nil, err
	}

	return book, nil
}

// FindMany retrieves a list of books from the database based on the provided request and status.
func (r *Repository) FindMany(request *api.Request, status string) (*Books, int, error) {
	books := &Books{}

	query := fmt.Sprintf(`
			SELECT b.* FROM books AS b 
		    WHERE b.deleted_at IS NULL %s LIMIT ? OFFSET ?
		`,
		r.statusQuery(status),
	)

	err := r.database.Select(books, query, request.Limit, request.Offset())
	if err != nil {
		return nil, 0, err
	}

	query = fmt.Sprintf(
		`SELECT COUNT(b.id) FROM books AS b WHERE b.deleted_at IS NULL %s`,
		r.statusQuery(status),
	)
	var count int
	err = r.database.Get(&count, query)
	if err != nil {
		return nil, 0, err
	}

	return books, count, nil
}

func (r *Repository) FindManyByName(req *api.Request, status string, search string) (*Books, int, error) {
	books := &Books{}
	query := fmt.Sprintf(`
			SELECT b.* FROM books AS b 
            WHERE b.name LIKE ? AND b.deleted_at IS NULL %s LIMIT ? OFFSET ?
		`,
		r.statusQuery(status),
	)

	err := r.database.Select(books, query, fmt.Sprintf("%%%s%%", search), req.Limit, req.Offset())
	if err != nil {
		return nil, 0, err
	}

	query = fmt.Sprintf(
		`SELECT COUNT(b.id) FROM books AS b WHERE b.name LIKE ? AND b.deleted_at IS NULL %s`,
		r.statusQuery(status),
	)
	var count int
	err = r.database.Get(&count, query, fmt.Sprintf("%%%s%%", search))
	if err != nil {
		return nil, 0, err
	}

	return books, count, nil
}

func (r *Repository) statusQuery(s string) string {
	switch s {
	case statusClosed:
		return " AND b.closed_at IS NOT NULL "
	case statusAny:
		return ""
	case statusActive:
		return " AND b.closed_at IS NULL "
	default:
		return " AND b.closed_at IS NULL "
	}
}

func (r *Repository) Create(book *Book) error {
	_, err := r.database.NamedExec(`
		INSERT INTO books (id, name, description, created_by, created_at)
		VALUES (:id, :name, :description, :created_by, :created_at)
	`, book)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Update(book *Book) error {
	_, err := r.database.NamedExec(`
		UPDATE books
		SET 
		    name = :name, 
		    description = :description,
			closed_at = :closed_at,
			deleted_at = :deleted_at
		WHERE id = :id
	`, book)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) AddObject(book *Book, m models.BookObject) (any, error) {
	e := buildObject(book, m)
	query := fmt.Sprintf(`
		INSERT INTO %s (id, book_id, %s, created_by, created_at)
		VALUES (NULL, :book_id, :%s, :created_by, :created_at)
	`, e.table(), e.field(), e.field())

	_, err := r.database.NamedExec(query, e)

	if err != nil {
		return nil, err
	}

	return e, nil
}

func CreateBook(d *sqlx.DB) *Repository {
	return &Repository{database: d}
}
