package model

import (
	"strings"

	"github.com/google/uuid"
)

type (
	Book struct {
		OrganizationId *uuid.UUID `validate:"required" json:"organization_id"`
		CreatedBy      *uuid.UUID
		Name           string  `validate:"required,gte=3" json:"name"`
		Description    *string `json:"description"`
	}

	BookObject interface {
		SetCreatedBy(id *uuid.UUID)
	}

	BookCategory struct {
		CategoryId *uuid.UUID `validate:"required" json:"category_id"`
		CreatedBy  *uuid.UUID
	}

	BookLocation struct {
		LocationId *uuid.UUID `validate:"required" json:"location_id"`
		CreatedBy  *uuid.UUID
	}
)

func (b *BookCategory) SetCreatedBy(id *uuid.UUID) {
	b.CreatedBy = id
	return
}

func (b *BookLocation) SetCreatedBy(id *uuid.UUID) {
	b.CreatedBy = id
	return
}

func DetermineBookObject(u string) BookObject {
	if strings.Contains(u, "categories") {
		return &BookCategory{}
	}

	return &BookLocation{}
}
