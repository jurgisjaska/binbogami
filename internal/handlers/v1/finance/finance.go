package finance

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
)

type Finance struct {
	echo     *echo.Group
	database *sqlx.DB
	auditlog *slog.Logger
}

func (h *Finance) initialize() *Finance {
	// repositories
	// routes

	return h
}
func CreateFinance(g *echo.Group, d *sqlx.DB, l *slog.Logger) *Finance {
	return (&Finance{echo: g, database: d, auditlog: l}).initialize()
}
