package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type (
	Organization struct {
		Id          *uuid.UUID `json:"id"`
		Name        *string    `json:"name"`
		Description *string    `json:"description"`
		CreatedBy   *uuid.UUID `json:"created_by"`
		CreatedAt   time.Time  `db:"created_at" json:"created_at"`
		UpdatedAt   *time.Time `db:"updated_at" json:"updated_at"`
		DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at"`
	}

	Organizations []Organization

	OrganizationRepository struct {
		database *sql.DB
	}
)
