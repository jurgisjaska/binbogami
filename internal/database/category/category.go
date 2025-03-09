package category

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api/models"
	"github.com/jurgisjaska/binbogami/internal/database/book"
)

type (
	Category struct {
		Id          *uuid.UUID `json:"id"`
		Name        string     `json:"name"`
		Description *string    `json:"description"`

		OrganizationId *uuid.UUID `db:"organization_id" json:"organization_id"`
		CreatedBy      *uuid.UUID `db:"created_by" json:"created_by"`

		CreatedAt time.Time  `db:"created_at" json:"created_at"`
		UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
		DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
	}

	Categories []Category

	// CategoryRepository is a struct that represents a repository for managing categories in the database.
	CategoryRepository struct {
		database *sqlx.DB
	}
)

// Find retrieves a category by its ID.
func (r *CategoryRepository) Find(id uuid.UUID) (*Category, error) {
	category := &Category{}
	err := r.database.Get(category, "SELECT * FROM categories WHERE id = ? AND deleted_at IS NULL", id.String())
	if err != nil {
		return nil, err
	}

	return category, nil
}

// ByBook retrieves a category by a book and category ID.
func (r *CategoryRepository) ByBook(book *book.Book, id *uuid.UUID) (*Category, error) {
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
func (r *CategoryRepository) ManyByBook(book *book.Book) (*Categories, error) {
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

// ManyByOrganization retrieves all categories for a given organization.
func (r *CategoryRepository) ManyByOrganization(org *uuid.UUID) (*Categories, error) {
	categories := &Categories{}
	query := `
		SELECT * 
		FROM categories 
		WHERE organization_id = ?
		AND deleted_at IS NULL
	`

	err := r.database.Select(categories, query, org)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoryRepository) Create(c *models.Category) (*Category, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	category := &Category{
		Id:             &id,
		Name:           c.Name,
		Description:    c.Description,
		OrganizationId: c.OrganizationId,
		CreatedBy:      c.CreatedBy,
		CreatedAt:      time.Now(),
	}

	_, err = r.database.NamedExec(`
		INSERT INTO categories (id, name, description, organization_id, created_by, created_at)
		VALUES (:id, :name, :description, :organization_id, :created_by, :created_at)
	`, category)

	if err != nil {
		return nil, err
	}

	return category, nil
}

func (r *CategoryRepository) Remove(category *Category) error {
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

// CreateCategory creates a new instance of CategoryRepository with the specified database connection.
func CreateCategory(d *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{database: d}
}
