package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	defaultOrganization int = iota + 1
)

type (
	Configuration struct {
		Id            *uuid.UUID `json:"id"`
		Configuration int        `json:"configuration"`
		Value         string     `json:"value"`

		CreatedBy *uuid.UUID `db:"created_by" json:"created_by"`

		CreatedAt time.Time  `db:"created_at" json:"createdAt"`
		UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt"`
	}

	// ConfigurationRepository represents a repository for storing user configuration data.
	ConfigurationRepository struct {
		database *sqlx.DB
	}
)

func (r *ConfigurationRepository) DefaultOrganization(user *User) (*Configuration, error) {
	return nil, nil
}

// CreateConfiguration creates a new instance of ConfigurationRepository with the specified SQL database connection.
func CreateConfiguration(d *sqlx.DB) *ConfigurationRepository {
	return &ConfigurationRepository{database: d}
}
