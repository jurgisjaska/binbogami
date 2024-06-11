package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	model "github.com/jurgisjaska/binbogami/internal/api/model/user"
)

const defaultOrganization int = iota + 1

type (
	Configuration struct {
		Id            *uuid.UUID `json:"id"`
		Configuration int        `json:"configuration"`
		Value         string     `json:"value"`

		UserId *uuid.UUID `db:"user_id" json:"userId"`

		CreatedAt time.Time  `db:"created_at" json:"createdAt"`
		UpdatedAt *time.Time `db:"updated_at" json:"updatedAt"`
	}

	// ConfigurationRepository represents a repository for storing user configuration data.
	ConfigurationRepository struct {
		database *sqlx.DB
	}
)

// FindDefaultOrganization retrieves the default organization configuration for a user.
func (r *ConfigurationRepository) FindDefaultOrganization(u *User) (*Configuration, error) {
	query := `
		SELECT * FROM user_configurations 
		WHERE configuration = ? AND user_id = ?
	`
	configuration := &Configuration{}
	err := r.database.Get(configuration, query, defaultOrganization, u.Id)
	if err != nil {
		return nil, err
	}

	return configuration, nil
}

// Upsert inserts a new configuration record into the database if it does not exist, or updates an existing record if it does.
func (r *ConfigurationRepository) Upsert(model *model.SetConfigurationRequest) (*Configuration, error) {
	id := uuid.New()
	configuration := &Configuration{
		Id:            &id,
		Configuration: model.Configuration,
		Value:         model.Value,
		UserId:        model.UserId,
		CreatedAt:     time.Now(),
	}

	_, err := r.database.NamedExec(`
		INSERT INTO user_configurations (id, configuration, value, user_id, created_at)
		VALUES (:id, :configuration, :value, :user_id, :created_at)
		ON DUPLICATE KEY UPDATE value = :value
	`, configuration)

	if err != nil {
		return nil, err
	}

	return configuration, nil
}

// CreateConfiguration creates a new instance of ConfigurationRepository.
func CreateConfiguration(d *sqlx.DB) *ConfigurationRepository {
	return &ConfigurationRepository{database: d}
}
