package models

import (
	"github.com/google/uuid"
)

type Entry struct {
	BookId     uuid.UUID  `validate:"required" json:"book_id"`
	CategoryId *uuid.UUID `validate:"required" json:"category_id"`
	LocationId *uuid.UUID `validate:"required" json:"location_id"`

	Amount      float64 `validate:"required,numeric" json:"amount"`
	Description *string `json:"description"`

	CreatedBy *uuid.UUID
}
