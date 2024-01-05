package internal

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func ConnectDatabase(c *Database) (*sqlx.DB, error) {
	return sqlx.Connect(
		"mysql",
		fmt.Sprintf(
			"%s:%s@(%s:%d)/%s?parseTime=true",
			c.Connection.Username,
			c.Connection.Password,
			c.Connection.Hostname,
			c.Connection.Port,
			c.Name,
		),
	)
}
