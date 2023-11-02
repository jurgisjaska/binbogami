package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type (
	Category struct {
		Id          *uuid.UUID `json:"id"`
		Name        *string    `json:"name"`
		Description *string    `json:"description"`
		CreatedAt   time.Time  `db:"created_at" json:"created_at"`
		UpdatedAt   *time.Time `db:"updated_at" json:"updated_at"`
		DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at"`
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

func (r *CategoryRepository) FindMany() (*Categories, error) {
	categories := &Categories{}
	err := r.database.Select(categories, "SELECT * FROM categories WHERE deleted_at IS NULL")
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoryRepository) Create(category *Category) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	category.Id = &id
	category.CreatedAt = time.Now()

	_, err = r.database.NamedExec(`
		INSERT INTO categories (id, name, description, created_at, updated_at, deleted_at)
		VALUES (:id, :name, :description, :created_at, :updated_at, :deleted_at)
	`, category)

	if err != nil {
		return err
	}

	return nil
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
