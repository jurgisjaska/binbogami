package v1

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/model"
	"github.com/jurgisjaska/binbogami/internal/database"
	"github.com/jurgisjaska/binbogami/internal/database/book"
	"github.com/jurgisjaska/binbogami/internal/database/member"
	"github.com/labstack/echo/v4"
)

type Category struct {
	echo       *echo.Group
	database   *sqlx.DB
	repository *database.CategoryRepository
	member     *member.MemberRepository
	book       *book.Repository
}

func (h *Category) initialize() *Category {
	h.repository = database.CreateCategory(h.database)
	h.member = member.CreateMember(h.database)
	h.book = book.CreateBook(h.database)

	h.echo.POST("/categories", h.create)
	h.echo.GET("/categories", h.byOrganization)
	h.echo.GET("/books/:id/categories", h.byBook)

	return h
}

func (h *Category) one(c echo.Context) error {
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

func (h *Category) byOrganization(c echo.Context) error {
	member, err := membership(h.member, c)
	if err != nil {
		return c.JSON(http.StatusForbidden, api.Error(err.Error()))
	}

	categories, err := h.repository.ManyByOrganization(member.OrganizationId)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("no categories found in the organization"))
	}

	return c.JSON(http.StatusOK, api.Success(categories, api.CreateRequest(c)))
}

func (h *Category) byBook(c echo.Context) error {
	_, err := membership(h.member, c)
	if err != nil {
		return c.JSON(http.StatusForbidden, api.Error(err.Error()))
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect book"))
	}

	book, err := h.book.Find(&id)
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
	member, err := membership(h.member, c)
	if err != nil {
		return c.JSON(http.StatusForbidden, api.Error(err.Error()))
	}

	category := &model.Category{}
	if err := c.Bind(category); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect category data"))
	}

	if err := c.Validate(category); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect category data", err.Error()))
	}

	category.CreatedBy = member.UserId
	category.OrganizationId = member.OrganizationId
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
func (h *Category) delete(c echo.Context) error {
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

func CreateCategory(g *echo.Group, d *sqlx.DB) *Category {
	return (&Category{echo: g, database: d}).initialize()
}
