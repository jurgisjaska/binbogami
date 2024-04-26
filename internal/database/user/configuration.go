package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	um "github.com/jurgisjaska/binbogami/internal/api/model/user"
)

const (
	defaultOrganization int = iota + 1
)

type (
	Configuration struct {
		Id            *uuid.UUID `json:"id"`
		Configuration int        `json:"configuration"`
		Value         string     `json:"value"`

		CreatedBy *uuid.UUID `db:"created_by" json:"createdBy"`

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
	configuration := &Configuration{}
	err := r.database.Get(configuration, `
		SELECT * FROM user_configurations 
		         WHERE configuration = ? AND created_by = ? AND deleted_at IS NULL
	`, defaultOrganization, user.Id)
	if err != nil {
		return nil, err
	}

	return configuration, nil
}

func (r *ConfigurationRepository) Create(model *um.SetConfiguration) (*Configuration, error) {
	id := uuid.New()
	configuration := &Configuration{
		Id:            &id,
		Configuration: model.Configuration,
		Value:         model.Value,
		CreatedBy:     model.CreatedBy,
		CreatedAt:     time.Now(),
	}

	_, err := r.database.NamedExec(`
		INSERT INTO user_configurations (id, configuration, value, created_by, created_at)
		VALUES (:id, :configuration, :value, :created_by, :created_at)
	`, configuration)

	if err != nil {
		return nil, err
	}

	return configuration, nil
}

// CreateConfiguration creates a new instance of ConfigurationRepository with the specified SQL database connection.
func CreateConfiguration(d *sqlx.DB) *ConfigurationRepository {
	return &ConfigurationRepository{database: d}
}
