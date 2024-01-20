package v1

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/model"
	"github.com/jurgisjaska/binbogami/internal/api/token"
	"github.com/jurgisjaska/binbogami/internal/database"
	"github.com/labstack/echo/v4"
)

type Category struct {
	echo       *echo.Group
	database   *sqlx.DB
	repository *database.CategoryRepository
	member     *database.MemberRepository
}

func (h *Category) initialize() *Category {
	h.repository = database.CreateCategory(h.database)
	h.member = database.CreateMember(h.database)

	h.echo.POST("/categories", h.create)

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

func (h *Category) many(c echo.Context) error {
	categories, err := h.repository.FindMany()
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("no categories found"))
	}

	return c.JSON(http.StatusOK, api.Success(categories, api.CreateRequest(c)))
}

func (h *Category) create(c echo.Context) error {
	category := &model.Category{}
	if err := c.Bind(category); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect category data"))
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(category); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect category data", err.Error()))
	}

	claims := token.FromContext(c)
	if claims.Id == nil {
		return c.JSON(http.StatusBadRequest, api.Error("invalid authentication token"))
	}

	member, err := h.member.Find(category.OrganizationId, claims.Id)
	if err != nil || member == nil {
		return c.JSON(http.StatusForbidden, api.Error("only organization members can create categories"))
	}

	category.CreatedBy = claims.Id
	entity, err := h.repository.Create(category)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(entity, api.CreateRequest(c)))
}

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
