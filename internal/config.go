package internal

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type (
	// Config represents the configuration for the application.
	// It contains environment variables, database configuration, and mail configuration.
	Config struct {
		Secret   string
		App      *URI
		Web      *URI
		Database *Database
		Mail     *Mail
		Auth     *Auth
		Loki     *Connection
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

	// URI represents the Uniform Resource Identifier configuration.
	// It holds the scheme, hostname, and port for a service.
	URI struct {
		Scheme   string
		Hostname string
		Port     int
	}

	// Auth represents the authentication configuration, including a secret and a URI for the authentication service.
	Auth struct {
		Secret string
		*URI
	}
)

// CreateConfig loads the configuration from the environment and creates an instance of config.
func CreateConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		Secret: os.Getenv("APP_SECRET"),
		App:    uri("APP"),
		Auth: &Auth{
			Secret: os.Getenv("AUTH_SERVICE_SECRET"),
			URI:    uri("AUTH_SERVICE"),
		},
		Web: uri("WEB_APPLICATION"),
		Database: &Database{
			Name: os.Getenv("DATABASE_NAME"),
			Connection: &Connection{
				Hostname: os.Getenv("DATABASE_HOSTNAME"),
				Port:     port(os.Getenv("DATABASE_PORT")),
				Username: os.Getenv("DATABASE_USERNAME"),
				Password: os.Getenv("DATABASE_PASSWORD"),
			},
		},
		Mail: &Mail{
			Sender: os.Getenv("MAIL_SENDER"),
			Connection: &Connection{
				Hostname: os.Getenv("MAIL_HOSTNAME"),
				Port:     port(os.Getenv("MAIL_PORT")),
				Username: os.Getenv("MAIL_USERNAME"),
				Password: os.Getenv("MAIL_PASSWORD"),
			},
		},
		Loki: &Connection{
			Hostname: os.Getenv("LOKI_HOSTNAME"),
			Port:     port(os.Getenv("LOKI_PORT")),
			Username: os.Getenv("LOKI_USERNAME"),
			Password: os.Getenv("LOKI_PASSWORD"),
		},
	}, nil
}

// port is a helper method that converts a string to an integer.
// It returns 0 if the string is empty or cannot be converted.
func port(s string) int {
	if s == "" {
		return 0
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return i
}

// uri is a helper method that creates a URI from environment variables with a given prefix.
func uri(prefix string) *URI {
	return &URI{
		Scheme:   os.Getenv(prefix + "_SCHEME"),
		Hostname: os.Getenv(prefix + "_HOSTNAME"),
		Port:     port(os.Getenv(prefix + "_PORT")),
	}
}
