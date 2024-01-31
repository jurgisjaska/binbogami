package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api/model"
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

	CategoryRepository struct {
		database *sqlx.DB
	}
)

func (r *CategoryRepository) Find(id uuid.UUID) (*Category, error) {
	category := &Category{}
	err := r.database.Get(category, "SELECT * FROM categories WHERE id = ? AND deleted_at IS NULL", id.String())
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (r *CategoryRepository) ByBook(book *Book, id *uuid.UUID) (*Category, error) {
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

func (r *CategoryRepository) FindMany() (*Categories, error) {
	categories := &Categories{}
	err := r.database.Select(categories, "SELECT * FROM categories WHERE deleted_at IS NULL")
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoryRepository) Create(c *model.Category) (*Category, error) {
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

	// @todo remove updated at and deleted at they are unnecessary as it will never be deleted on creation
	_, err = r.database.NamedExec(`
		INSERT INTO categories (id, name, description, organization_id, created_by, created_at, updated_at, deleted_at)
		VALUES (:id, :name, :description, :organization_id, :created_by, :created_at, :updated_at, :deleted_at)
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

func CreateCategory(d *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{database: d}
}
