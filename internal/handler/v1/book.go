package v1

import (
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/database"
	"github.com/labstack/echo/v4"
)

// Book represents a book handler.
type Book struct {
	echo       *echo.Group
	database   *sqlx.DB
	repository *database.BookRepository
}

func (h *Book) initialize() *Book {
	h.repository = database.CreateBook(h.database)
	// h.echo.GET("/books/:id", h.one)
	// h.echo.GET("/books", h.many)
	// h.echo.POST("/books", h.create)
	// h.echo.PUT("/books/:id", h.update)
	// h.echo.DELETE("/books/:id", h.delete)

	return h
}

// CreateBook creates a new instance of Book handler.
func CreateBook(g *echo.Group, d *sqlx.DB) *Book {
	return (&Book{echo: g, database: d}).initialize()
}
