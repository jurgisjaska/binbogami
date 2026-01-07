package configuration

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	um "github.com/jurgisjaska/binbogami/internal/api/models/user"
)

// ConfigurationRepository represents a repository for storing user configuration data.
type ConfigurationRepository struct {
	database *sqlx.DB
}

// Upsert inserts a new configuration record into the database if it does not exist, or updates an existing record if it does.
func (r *ConfigurationRepository) Upsert(model *um.SetConfigurationRequest) (*Configuration, error) {
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
