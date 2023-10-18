package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type (
	Category struct {
		Id          *uuid.UUID `json:"id"`
		Name        string     `json:"name"`
		Description string     `json:"description"`
		CreatedAt   time.Time  `db:"created_at" json:"created_at"`
		UpdatedAt   *time.Time `db:"updated_at" json:"updated_at"`
		DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at"`
	}

	Categories []*Category

	CategoryRepository struct {
		database *sqlx.DB
	}
)

func (r *CategoryRepository) Find(id uuid.UUID) (*Category, error) {
	category := &Category{}
	err := r.database.Get(category, "SELECT * FROM categories WHERE id = ?", id.String())
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (r *CategoryRepository) FindMany() *Categories {
	return nil
}

func CreateCategory(d *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{database: d}
}
