package models

import (
	"strings"

	"github.com/google/uuid"
)

type (
	CreateBook struct {
		Name        string  `validate:"required,gte=3" json:"name"`
		Description *string `json:"description"`
	}

	UpdateBook struct {
		Id          uuid.UUID `validate:"required,uuid" json:"id"`
		Name        string    `validate:"required,gte=3" json:"name"`
		Description *string   `json:"description"`
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
