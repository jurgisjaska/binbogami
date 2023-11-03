package handler

import (
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/app/database"
	"github.com/labstack/echo/v4"
)

type (
	Book struct {
		echo       *echo.Echo
		database   *sqlx.DB
		repository *database.BookRepository
	}
)

func (h *Book) initialize() *Book {
	h.repository = database.CreateBook(h.database)
	// h.echo.GET("/books/:id", h.one)
	// h.echo.GET("/books", h.many)
	// h.echo.POST("/books", h.create)
	// h.echo.PUT("/books/:id", h.update)
	// h.echo.DELETE("/books/:id", h.delete)

	return h
}

func CreateBook(e *echo.Echo, d *sqlx.DB) *Book {
	return (&Book{echo: e, database: d}).initialize()
}
