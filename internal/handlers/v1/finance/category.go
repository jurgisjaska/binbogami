package finance

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/models"
	"github.com/jurgisjaska/binbogami/internal/database/category"
	"github.com/jurgisjaska/binbogami/internal/database/entry"
	"github.com/labstack/echo/v5"
)

type Category struct {
	echo            *echo.Group
	database        *sqlx.DB
	repository      category.CategoryRepository
	entryRepository entry.EntryRepository
	auditlog        *slog.Logger
}

// initialize sets up routes and dependencies for the Category handler and returns the initialized handler instance.
func (h *Category) initialize() *Category {
	h.repository = category.CreateCategory(h.database)
	h.entryRepository = entry.CreateEntry(h.database)

	h.echo.GET("/categories", h.index)
	// stats: this month + change from last month / monthly graph for a year
	h.echo.GET("/categories/:id", h.show)

	h.echo.POST("/categories", h.create)
	h.echo.PUT("/categories/:id", h.update)

	h.echo.DELETE("/categories/:id", h.destroy)

	h.echo.GET("/categories/:id/entries", h.entries)

	return h
}

func (h *Category) index(c *echo.Context) error {
	request := api.CreateRequest(c)

	var categories *category.Categories
	var t int
	var err error

	categories, t, err = h.repository.FindMany(request)
	if err != nil {
		h.auditlog.Warn("category index error: failed to find categories", "error", err.Error())
		return c.JSON(http.StatusNotFound, api.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(categories, request, t))
}

// show retrieves a category by ID, fetches the category from the database, and returns a JSON response.
func (h *Category) show(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.auditlog.Warn("category show error: incorrect category", "error", err.Error())
		return c.JSON(http.StatusBadRequest, api.Error("incorrect category"))
	}

	entity, err := h.repository.Find(id)
	if err != nil {
		h.auditlog.Warn("category show error: category not found", "error", err.Error(), "category_id", id)
		return c.JSON(http.StatusNotFound, api.Error("category not found"))
	}

	return c.JSON(http.StatusOK, api.Success(entity, api.CreateRequest(c)))
}

func (h *Category) entries(c *echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.auditlog.Warn("category entries error: incorrect category", "error", err.Error())
		return c.JSON(http.StatusBadRequest, api.Error("incorrect category"))
	}

	entity, err := h.repository.Find(id)
	if err != nil {
		h.auditlog.Warn("category entries error: category not found", "error", err.Error(), "category_id", id)
		return c.JSON(http.StatusNotFound, api.Error("category not found"))
	}

	request := api.CreateRequest(c)

	var entries *entry.Entries
	var t int

	entries, t, err = h.entryRepository.FindManyByCategory(entity.Id, request)
	if err != nil {
		h.auditlog.Warn("category entries error: failed to find entries", "error", err.Error())
		return c.JSON(http.StatusNotFound, api.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(entries, request, t))
}

func (h *Category) create(c *echo.Context) error {
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
func (h *Category) update(c *echo.Context) error {
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
func (h *Category) destroy(c *echo.Context) error {
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

// CreateCategory initializes the Category resource, sets up its repository dependencies, and maps HTTP endpoints.
func CreateCategory(g *echo.Group, d *sqlx.DB, l *slog.Logger) *Category {
	return (&Category{echo: g, database: d, auditlog: l}).initialize()
}
