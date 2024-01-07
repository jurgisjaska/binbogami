package internal

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		Environment string
		Secret      string
		Hostname    string
		Port        int

		Database *Database
		Mail     *Mail
	}

	// Database represents the database configuration.
	// It holds the database name and connection configuration.
	Database struct {
		Name       string
		Connection *Connection
	}

	// Mail represents the SMTP mail client configuration.
	// It holds the connection information and the sender email address.
	Mail struct {
		Sender     string
		Connection *Connection
	}

	// Connection represents the connection configuration for local or 3rd party service.
	Connection struct {
		Hostname string
		Port     int
		Username string
		Password string
	}
)

// CreateConfig loads the configuration from the environment and creates an instance of config.
func CreateConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	p, _ := strconv.Atoi(os.Getenv("DATABASE_PORT"))
	ap, _ := strconv.Atoi(os.Getenv("APP_PORT"))
	mp, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))

	return &Config{
		Environment: os.Getenv("APP_ENVIRONMENT"),
		Secret:      os.Getenv("APP_SECRET"),
		Hostname:    os.Getenv("APP_HOSTNAME"),
		Port:        ap,
		Database: &Database{
			Name: os.Getenv("DATABASE_NAME"),
			Connection: &Connection{
				Hostname: os.Getenv("DATABASE_HOSTNAME"),
				Port:     p,
				Username: os.Getenv("DATABASE_USERNAME"),
				Password: os.Getenv("DATABASE_PASSWORD"),
			},
		},
		Mail: &Mail{
			Sender: os.Getenv("MAIL_SENDER"),
			Connection: &Connection{
				Hostname: os.Getenv("MAIL_HOSTNAME"),
				Port:     mp,
				Username: os.Getenv("MAIL_USERNAME"),
				Password: os.Getenv("MAIL_PASSWORD"),
			},
		},
	}, nil
}
