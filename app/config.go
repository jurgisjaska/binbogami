package app

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		Environment string
		Salt        string

		Database *Database
	}

	Database struct {
		Name       string
		Connection *Connection
	}

	Connection struct {
		Hostname string
		Port     int
		Username string
		Password string
	}
)

func CreateConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	p, _ := strconv.Atoi(os.Getenv("DATABASE_PORT"))

	return &Config{
		Environment: os.Getenv("APP_ENVIRONMENT"),
		Salt:        os.Getenv("APP_SALT"),
		Database: &Database{
			Name: os.Getenv("DATABASE_NAME"),
			Connection: &Connection{
				Hostname: os.Getenv("DATABASE_HOSTNAME"),
				Port:     p,
				Username: os.Getenv("DATABASE_USERNAME"),
				Password: os.Getenv("DATABASE_PASSWORD"),
			},
		},
	}, nil
}
