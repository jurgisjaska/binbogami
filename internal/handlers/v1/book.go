package v1

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/models"
	"github.com/jurgisjaska/binbogami/internal/database/book"
	"github.com/labstack/echo/v4"
)

// Book represents a book handlers.
type Book struct {
	echo       *echo.Group
	database   *sqlx.DB
	repository *book.Repository
}

func (h *Book) initialize() *Book {
	h.repository = book.CreateBook(h.database)

	h.echo.GET("/books", h.index)
	h.echo.POST("/books", h.create)
	h.echo.PUT("/books", h.update)
	h.echo.GET("/books/:id", h.show)
	h.echo.POST("/books/:id/categories", h.add)
	h.echo.POST("/books/:id/locations", h.add)

	return h
}

func (h *Book) index(c echo.Context) error {
	req := api.CreateRequest(c)
	status := c.QueryParam("status")

	books, t, err := h.repository.FindMany(req, status)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("no books found"))
	}

	return c.JSON(http.StatusOK, api.Success(books, req, t))
}

func (h *Book) create(c echo.Context) error {
	bm := &models.CreateBook{}
	// bm.CreatedBy = member.UserId

	if err := c.Bind(bm); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect book data"))
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(bm); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect book data", err.Error()))
	}

	entity, err := h.repository.Create(bm)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(entity, api.CreateRequest(c)))
}

func (h *Book) update(c echo.Context) error {
	bm := &models.UpdateBook{}
	if err := c.Bind(bm); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect book data"))
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(bm); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect book data", err.Error()))
	}

	b, err := h.repository.Find(bm.Id)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error(err.Error()))
	}

	entity, err := h.repository.Update(b, bm)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(entity, api.CreateRequest(c)))
}

func (h *Book) add(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect book"))
	}

	book, err := h.repository.Find(&id)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("book not found"))
	}

	m := models.DetermineBookObject(c.Request().URL.String())
	if err := c.Bind(m); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect book data"))
	}

	if err := c.Validate(m); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect book data", err.Error()))
	}

	// m.SetCreatedBy(member.UserId)
	entity, err := h.repository.AddObject(book, m)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(entity, api.CreateRequest(c)))
}

func (h *Book) show(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect book"))
	}

	entity, err := h.repository.Find(&id)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("no books found"))
	}

	return c.JSON(http.StatusOK, api.Success(entity, api.CreateRequest(c)))
}

// CreateBook creates a new instance of Book handlers.
func CreateBook(g *echo.Group, d *sqlx.DB) *Book {
	return (&Book{echo: g, database: d}).initialize()
}
