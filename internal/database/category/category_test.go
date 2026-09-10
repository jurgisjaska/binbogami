package category

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCategoryStruct(t *testing.T) {
	id := uuid.New()
	createdBy := uuid.MustParse("05e7257a-b21c-11ee-9a7a-5ab75f0c1cab")
	now := time.Now()
	desc := "P90 ammo, C4, and standard issue gear"
	color := "#4CAF50"
	icon := "boxes-stacked"

	tests := []struct {
		name     string
		category *Category
	}{
		{
			name: "Category with description, color, and icon",
			category: &Category{
				Id:          &id,
				Name:        "Mission Supplies",
				Description: &desc,
				Color:       &color,
				Icon:        &icon,
				CreatedBy:   &createdBy,
				CreatedAt:   now,
				UpdatedAt:   nil,
				DeletedAt:   nil,
			},
		},
		{
			name: "Category without optional description, color, and icon",
			category: &Category{
				Id:          &id,
				Name:        "Off-World Recon",
				Description: nil,
				Color:       nil,
				Icon:        nil,
				CreatedBy:   &createdBy,
				CreatedAt:   now,
				UpdatedAt:   nil,
				DeletedAt:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.category.Id)
			assert.NotEmpty(t, tt.category.Name)
			assert.Equal(t, &createdBy, tt.category.CreatedBy)
			if tt.category.Description != nil {
				assert.Equal(t, desc, *tt.category.Description)
			}
			if tt.category.Color != nil {
				assert.Equal(t, color, *tt.category.Color)
			}
			if tt.category.Icon != nil {
				assert.Equal(t, icon, *tt.category.Icon)
			}
		})
	}
}

func TestSortable(t *testing.T) {
	expectedFields := []string{"name", "description", "created_at"}
	for _, field := range expectedFields {
		t.Run("Field_"+field+"_is_sortable", func(t *testing.T) {
			assert.True(t, sortable[field])
		})
	}

	unexpectedFields := []string{"id", "color", "icon", "created_by", "deleted_at", "updated_at", "invalid"}
	for _, field := range unexpectedFields {
		t.Run("Field_"+field+"_is_not_sortable", func(t *testing.T) {
			assert.False(t, sortable[field])
		})
	}
}
