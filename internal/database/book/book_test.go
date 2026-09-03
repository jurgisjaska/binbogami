package book

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jurgisjaska/binbogami/internal/api/models"
	"github.com/stretchr/testify/assert"
)

func TestBookStruct(t *testing.T) {
	id := uuid.New()
	createdBy := uuid.MustParse("05e7257a-b21c-11ee-9a7a-5ab75f0c1cab")
	now := time.Now()
	desc := "The primary financial ledger for Stargate Command operations"
	author := "Jack O'Neil"

	tests := []struct {
		name string
		book *Book
	}{
		{
			name: "Active Book with description and author",
			book: &Book{
				Id:          id,
				Name:        "Year 2026",
				Description: &desc,
				CreatedBy:   createdBy,
				Author:      &author,
				CreatedAt:   now,
				UpdatedAt:   nil,
				DeletedAt:   nil,
				ClosedAt:    nil,
			},
		},
		{
			name: "Closed Book with closed timestamp",
			book: &Book{
				Id:          id,
				Name:        "Year 2025",
				Description: &desc,
				CreatedBy:   createdBy,
				Author:      nil,
				CreatedAt:   now,
				UpdatedAt:   nil,
				DeletedAt:   nil,
				ClosedAt:    &now,
			},
		},
		{
			name: "Soft-deleted Book without optional description",
			book: &Book{
				Id:          id,
				Name:        "Year 2025 deleted",
				Description: nil,
				CreatedBy:   createdBy,
				CreatedAt:   now,
				UpdatedAt:   nil,
				DeletedAt:   &now,
				ClosedAt:    nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.book.Id)
			assert.NotEmpty(t, tt.book.Name)
			assert.Equal(t, createdBy, tt.book.CreatedBy)
			if tt.book.Description != nil {
				assert.Equal(t, desc, *tt.book.Description)
			}
			if tt.book.ClosedAt != nil {
				assert.Equal(t, now, *tt.book.ClosedAt)
			}
			if tt.book.DeletedAt != nil {
				assert.Equal(t, now, *tt.book.DeletedAt)
			}
		})
	}
}

func TestStatusQuery(t *testing.T) {
	repo := &Repository{}

	tests := []struct {
		name     string
		status   string
		expected string
	}{
		{
			name:     "Status active returns closed_at IS NULL condition",
			status:   statusActive,
			expected: " AND b.closed_at IS NULL ",
		},
		{
			name:     "Status closed returns closed_at IS NOT NULL condition",
			status:   statusClosed,
			expected: " AND b.closed_at IS NOT NULL ",
		},
		{
			name:     "Status any returns empty condition",
			status:   statusAny,
			expected: "",
		},
		{
			name:     "Unknown status defaults to active condition",
			status:   "unknown_status",
			expected: " AND b.closed_at IS NULL ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.statusQuery(tt.status)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildObject(t *testing.T) {
	b := &Book{
		Id: uuid.New(),
	}
	catID := uuid.New()
	locID := uuid.New()
	userID := uuid.New()

	t.Run("Build BookCategory object", func(t *testing.T) {
		catModel := &models.BookCategory{
			CategoryId: &catID,
			CreatedBy:  &userID,
		}
		obj := buildObject(b, catModel)
		assert.NotNil(t, obj)
		assert.Equal(t, "books_categories", obj.table())
		assert.Equal(t, "category_id", obj.field())
	})

	t.Run("Build BookLocation object", func(t *testing.T) {
		locModel := &models.BookLocation{
			LocationId: &locID,
			CreatedBy:  &userID,
		}
		obj := buildObject(b, locModel)
		assert.NotNil(t, obj)
		assert.Equal(t, "books_locations", obj.table())
		assert.Equal(t, "location_id", obj.field())
	})
}
