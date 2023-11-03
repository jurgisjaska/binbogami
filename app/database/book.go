package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type (
	Book struct {
		Id          *uuid.UUID `json:"id"`
		Name        *string    `json:"name"`
		Description *string    `json:"description"`
		// owned_by - organization
		CreatedAt time.Time  `db:"created_at" json:"created_at"`
		UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
		DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
	}

	Books []Book

	BookRepository struct {
		database *sqlx.DB
	}
)

func CreateBook(d *sqlx.DB) *BookRepository {
	return &BookRepository{database: d}
}
