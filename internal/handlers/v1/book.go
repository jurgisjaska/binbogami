package v1

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/models"
	"github.com/jurgisjaska/binbogami/internal/database/book"
	"github.com/jurgisjaska/binbogami/internal/database/user"
	"github.com/labstack/echo/v4"
)

const (
	destroyDelete = "delete"
	destroyClose  = "close"
)

// Book represents a book handlers.
type Book struct {
	echo           *echo.Group
	database       *sqlx.DB
	repository     *book.Repository
	userRepository *user.Repository
}

func (h *Book) initialize() *Book {
	h.repository = book.CreateBook(h.database)
	h.userRepository = user.CreateUser(h.database)

	h.echo.GET("/books", h.index)
	h.echo.POST("/books", h.create)
	h.echo.PUT("/books/:id", h.update)
	h.echo.GET("/books/:id", h.show)
	h.echo.DELETE("/books/:id", h.destroy)
	h.echo.POST("/books/:id/categories", h.add)
	h.echo.POST("/books/:id/locations", h.add)

	return h
}

func (h *Book) index(c echo.Context) error {
	req := api.CreateRequest(c)
	status := c.QueryParam("status")
	query := c.QueryParam("query")

	var books *book.Books
	var t int
	var err error

	if len(query) == 0 {
		books, t, err = h.repository.FindMany(req, status)
	} else {
		books, t, err = h.repository.FindManyByName(req, status, query)
	}

	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("no books found"))
	}

	return c.JSON(http.StatusOK, api.Success(books, req, t))
}

func (h *Book) create(c echo.Context) error {
	request := &models.CreateBook{}
	u, err := currentUser(h.userRepository, c)
	if err != nil {
		return c.JSON(http.StatusForbidden, api.Error(err.Error()))
	}

	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect book data"))
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(request); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect book data", err.Error()))
	}

	book := &book.Book{
		Id:          uuid.New(),
		Name:        request.Name,
		Description: request.Description,
		CreatedBy:   u.Id,
		CreatedAt:   time.Now(),
	}

	err = h.repository.Create(book)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(book, api.CreateRequest(c)))
}

func (h *Book) update(c echo.Context) error {
	request := &models.UpdateBook{}
	_, err := currentUser(h.userRepository, c)
	if err != nil {
		return c.JSON(http.StatusForbidden, api.Error(err.Error()))
	}

	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect book data"))
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(request); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect book data", err.Error()))
	}

	book, err := h.repository.Find(request.Id)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error(err.Error()))
	}

	book.Name = request.Name
	book.Description = request.Description

	err = h.repository.Update(book)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(book, api.CreateRequest(c)))
}

func (h *Book) destroy(c echo.Context) error {
	_, err := currentUser(h.userRepository, c)
	if err != nil {
		return c.JSON(http.StatusForbidden, api.Error(err.Error()))
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect book"))
	}

	entity, err := h.repository.Find(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("no books found"))
	}

	n := time.Now()
	t := c.QueryParam("type")

	switch t {
	case destroyClose:
		entity.ClosedAt = &n
	case destroyDelete:
		entity.DeletedAt = &n
	default:
		entity.DeletedAt = &n
	}

	err = h.repository.Update(entity)
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

	book, err := h.repository.Find(id)
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

	entity, err := h.repository.Find(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("no books found"))
	}

	return c.JSON(http.StatusOK, api.Success(entity, api.CreateRequest(c)))
}

// CreateBook creates a new instance of Book handlers.
func CreateBook(g *echo.Group, d *sqlx.DB) *Book {
	return (&Book{echo: g, database: d}).initialize()
}
