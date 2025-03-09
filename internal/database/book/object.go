package book

import (
	"time"

	"github.com/google/uuid"
	"github.com/jurgisjaska/binbogami/internal/api/models"
)

type (
	object interface {
		table() string
		field() string
	}

	category struct {
		Id         int        `json:"id"`
		BookId     *uuid.UUID `db:"book_id" json:"book_id"`
		CategoryId *uuid.UUID `db:"category_id" json:"category_id"`

		CreatedBy *uuid.UUID `db:"created_by" json:"created_by"`

		CreatedAt time.Time  `db:"created_at" json:"created_at"`
		DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
	}

	location struct {
		Id         int        `json:"id"`
		BookId     *uuid.UUID `db:"book_id" json:"book_id"`
		LocationId *uuid.UUID `db:"location_id" json:"location_id"`

		CreatedBy *uuid.UUID `db:"created_by" json:"created_by"`

		CreatedAt time.Time  `db:"created_at" json:"created_at"`
		DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
	}
)

func buildObject(book *Book, m models.BookObject) object {
	if obj, ok := m.(*models.BookCategory); ok {
		return category{
			BookId:     book.Id,
			CategoryId: obj.CategoryId,
			CreatedBy:  obj.CreatedBy,
			CreatedAt:  time.Now(),
		}
	} else if obj, ok := m.(*models.BookLocation); ok {
		return location{
			BookId:     book.Id,
			LocationId: obj.LocationId,
			CreatedBy:  obj.CreatedBy,
			CreatedAt:  time.Now(),
		}
	}

	return nil
}

func (b category) table() string {
	return "books_categories"
}

func (b location) table() string {
	return "books_locations"
}

func (b category) field() string {
	return "category_id"
}

func (b location) field() string {
	return "location_id"
}
