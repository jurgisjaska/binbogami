package v1

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/models"
	"github.com/jurgisjaska/binbogami/internal/database/book"
	"github.com/jurgisjaska/binbogami/internal/database/category"
	"github.com/labstack/echo/v4"
)

type Category struct {
	echo       *echo.Group
	database   *sqlx.DB
	repository *category.Repository
	book       *book.Repository
}

// initialize sets up routes and dependencies for the Category handler and returns the initialized handler instance.
func (h *Category) initialize() *Category {
	h.repository = category.CreateCategory(h.database)
	h.book = book.CreateBook(h.database)

	h.echo.GET("/categories", h.index)
	h.echo.POST("/categories", h.create)
	h.echo.PUT("/categories/:id", h.update)
	h.echo.GET("/categories/:id", h.show)
	h.echo.DELETE("/categories/:id", h.destroy)

	// @todo remove this and prepare propper endpoint.
	h.echo.GET("/books/:id/categories", h.byBook)

	return h
}

func (h *Category) index(c echo.Context) error {
	request := api.CreateRequest(c)

	var categories *category.Categories
	var t int
	var err error

	categories, t, err = h.repository.FindMany(request)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(categories, request, t))
}

// show retrieves a category by ID, fetches the category from the database, and returns a JSON response.
func (h *Category) show(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect category id"))
	}

	entity, err := h.repository.Find(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("category not found"))
	}

	return c.JSON(http.StatusOK, api.Success(entity, api.CreateRequest(c)))
}

// @deprecated
// @todo remove
func (h *Category) byBook(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect book"))
	}

	book, err := h.book.Find(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("book not found"))
	}

	categories, err := h.repository.ManyByBook(book)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("no categories found"))
	}

	return c.JSON(http.StatusOK, api.Success(categories, api.CreateRequest(c)))
}

func (h *Category) create(c echo.Context) error {
	category := &models.Category{}
	if err := c.Bind(category); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect category data"))
	}

	if err := c.Validate(category); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect category data", err.Error()))
	}

	// category.CreatedBy = member.UserId
	entity, err := h.repository.Create(category)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(entity, api.CreateRequest(c)))
}

// @deprecated
func (h *Category) update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect category"))
	}

	category, err := h.repository.Find(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("category not found"))
	}

	return c.JSON(http.StatusOK, api.Success(category, api.CreateRequest(c)))
}

// @deprecated
func (h *Category) destroy(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect category"))
	}

	category, err := h.repository.Find(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("category not found"))
	}

	if err = h.repository.Remove(category); err != nil {
		return c.JSON(http.StatusInternalServerError, api.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(true, api.CreateRequest(c)))
}

// CreateCategory initializes the Category resource, sets up its repository dependencies, and maps its HTTP endpoints.
func CreateCategory(g *echo.Group, d *sqlx.DB) *Category {
	return (&Category{echo: g, database: d}).initialize()
}
