package location

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestLocationStruct(t *testing.T) {
	id := uuid.New()
	createdBy := uuid.MustParse("05e7257a-b21c-11ee-9a7a-5ab75f0c1cab")
	now := time.Now()
	desc := "SGC Headquarters Base"
	addr := "1 Norad Rd, Colorado Springs, CO 80906, USA"

	tests := []struct {
		name     string
		location *Location
	}{
		{
			name: "Location with description and address",
			location: &Location{
				Id:          &id,
				Name:        "Cheyenne Mountain Complex",
				Description: &desc,
				Address:     &addr,
				CreatedBy:   &createdBy,
				CreatedAt:   now,
				UpdatedAt:   nil,
				DeletedAt:   nil,
			},
		},
		{
			name: "Location without optional description and address",
			location: &Location{
				Id:          &id,
				Name:        "Alpha Site",
				Description: nil,
				Address:     nil,
				CreatedBy:   &createdBy,
				CreatedAt:   now,
				UpdatedAt:   nil,
				DeletedAt:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.location.Id)
			assert.NotEmpty(t, tt.location.Name)
			assert.Equal(t, &createdBy, tt.location.CreatedBy)
			if tt.location.Address != nil {
				assert.Equal(t, addr, *tt.location.Address)
			}
		})
	}
}
