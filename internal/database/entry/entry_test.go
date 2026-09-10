package entry

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEntryStruct(t *testing.T) {
	id := uuid.New()
	bookID := uuid.MustParse("7b3a1f90-2c4d-4e5f-8a1b-9c0d1e2f3a4b")
	catID := uuid.MustParse("f9e8d7c6-b5a4-4321-80f1-e2d3c4b5a697")
	locID := uuid.MustParse("e0000000-0000-4000-8000-000000000002")
	createdBy := uuid.MustParse("05e7257a-b21c-11ee-9a7a-5ab75f0c1cab")
	now := time.Now()
	desc := "P90 5.7x28mm ammunition crates for SG-1 expedition"

	tests := []struct {
		name  string
		entry *Entry
	}{
		{
			name: "Entry with description and optional fields",
			entry: &Entry{
				Id:          &id,
				Amount:      9600.42,
				Description: &desc,
				BookId:      bookID,
				CategoryId:  &catID,
				LocationId:  &locID,
				CreatedBy:   &createdBy,
				CreatedAt:   now,
				UpdatedAt:   nil,
				DeletedAt:   nil,
			},
		},
		{
			name: "Entry without optional description",
			entry: &Entry{
				Id:          &id,
				Amount:      150.00,
				Description: nil,
				BookId:      bookID,
				CategoryId:  &catID,
				LocationId:  &locID,
				CreatedBy:   &createdBy,
				CreatedAt:   now,
				UpdatedAt:   nil,
				DeletedAt:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.entry.Id)
			assert.Greater(t, tt.entry.Amount, float64(0))
			assert.Equal(t, bookID, tt.entry.BookId)
			assert.Equal(t, &createdBy, tt.entry.CreatedBy)
			if tt.entry.Description != nil {
				assert.Equal(t, desc, *tt.entry.Description)
			}
		})
	}
}

func TestSortable(t *testing.T) {
	expectedFields := []string{"amount", "description", "created_at"}
	for _, field := range expectedFields {
		t.Run("Field_"+field+"_is_sortable", func(t *testing.T) {
			assert.True(t, sortable[field])
		})
	}

	unexpectedFields := []string{"id", "book_id", "category_id", "location_id", "created_by", "deleted_at", "updated_at", "invalid"}
	for _, field := range unexpectedFields {
		t.Run("Field_"+field+"_is_not_sortable", func(t *testing.T) {
			assert.False(t, sortable[field])
		})
	}
}
