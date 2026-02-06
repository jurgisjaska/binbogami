package category

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/models"
	"github.com/jurgisjaska/binbogami/internal/database/book"
)

type (
	Category struct {
		Id          *uuid.UUID `json:"id"`
		Name        string     `json:"name"`
		Description *string    `json:"description"`
		Color       *string    `json:"color"`

		CreatedBy *uuid.UUID `db:"created_by" json:"created_by"`

		CreatedAt time.Time  `db:"created_at" json:"created_at"`
		UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
		DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
	}

	Categories []Category

	// Repository is a struct that represents a repository for managing categories in the database.
	Repository struct {
		database *sqlx.DB
	}
)

// Find retrieves a category by its ID.
func (r *Repository) Find(id uuid.UUID) (*Category, error) {
	category := &Category{}
	err := r.database.Get(category, "SELECT * FROM categories WHERE id = ? AND deleted_at IS NULL", id.String())
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (r *Repository) FindMany(request *api.Request) (*Categories, int, error) {
	categories := &Categories{}
	query := `SELECT * FROM categories WHERE deleted_at IS NULL LIMIT ? OFFSET ?`

	err := r.database.Select(categories, query, request.Limit, request.Offset())
	if err != nil {
		return nil, 0, err
	}

	query = `SELECT COUNT(id) FROM categories WHERE deleted_at IS NULL`
	var count int
	err = r.database.Get(&count, query)
	if err != nil {
		return nil, 0, err
	}

	return categories, count, nil
}

// ByBook retrieves a category by a book and category ID.
func (r *Repository) ByBook(book *book.Book, id *uuid.UUID) (*Category, error) {
	query := `
		SELECT categories.* 
		FROM categories 
		JOIN books_categories AS bc ON bc.category_id = categories.id
		JOIN books AS b ON b.id = bc.book_id
		WHERE 
		    b.id = ? AND categories.id = ?
		    AND categories.deleted_at IS NULL
			AND bc.deleted_at IS NULL AND b.deleted_at IS NULL
	`

	category := &Category{}
	if err := r.database.Get(category, query, book.Id, id); err != nil {
		return nil, err
	}

	return category, nil
}

// ManyByBook retrieves categories associated with a book.
func (r *Repository) ManyByBook(book *book.Book) (*Categories, error) {
	categories := &Categories{}
	query := `
		SELECT categories.* 
		FROM categories 
		JOIN books_categories AS bc ON bc.category_id = categories.id
		JOIN books AS b ON b.id = bc.book_id
		WHERE 
		    b.id = ?
		    AND categories.deleted_at IS NULL
			AND bc.deleted_at IS NULL AND b.deleted_at IS NULL
	`

	err := r.database.Select(categories, query, book.Id)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *Repository) Create(c *models.Category) (*Category, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	category := &Category{
		Id:          &id,
		Name:        c.Name,
		Description: c.Description,
		CreatedBy:   c.CreatedBy,
		CreatedAt:   time.Now(),
	}

	_, err = r.database.NamedExec(`
		INSERT INTO categories (id, name, description, created_by, created_at)
		VALUES (:id, :name, :description, :created_by, :created_at)
	`, category)

	if err != nil {
		return nil, err
	}

	return category, nil
}

func (r *Repository) Remove(category *Category) error {
	_, err := r.database.NamedExec(
		`UPDATE categories SET deleted_at = :deleted WHERE id = :id`,
		map[string]interface{}{
			"deleted": time.Now(),
			"id":      category.Id,
		})

	if err != nil {
		return err
	}

	return nil
}

// CreateCategory creates a new instance of Repository with the specified database connection.
func CreateCategory(d *sqlx.DB) *Repository {
	return &Repository{database: d}
}
