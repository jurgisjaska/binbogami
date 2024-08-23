package v1

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/model"
	"github.com/jurgisjaska/binbogami/internal/database"
	"github.com/jurgisjaska/binbogami/internal/database/book"
	"github.com/labstack/echo/v4"
)

// Book represents a book handler.
type Book struct {
	echo         *echo.Group
	database     *sqlx.DB
	repository   *book.Repository
	organization *database.OrganizationRepository
	member       *database.MemberRepository
}

func (h *Book) initialize() *Book {
	h.repository = book.CreateBook(h.database)
	h.organization = database.CreateOrganization(h.database)
	h.member = database.CreateMember(h.database)

	h.echo.POST("/books", h.create)
	h.echo.GET("/books", h.byOrganization)
	h.echo.POST("/books/:id/categories", h.add)
	h.echo.POST("/books/:id/locations", h.add)

	return h
}

func (h *Book) create(c echo.Context) error {
	member, err := membership(h.member, c)
	if err != nil {
		return c.JSON(http.StatusForbidden, api.Error(err.Error()))
	}

	book := &model.Book{}
	if err := c.Bind(book); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect book data"))
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(book); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect book data", err.Error()))
	}

	book.CreatedBy = member.UserId
	book.OrganizationId = member.OrganizationId
	entity, err := h.repository.Create(book)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(entity, api.CreateRequest(c)))
}

func (h *Book) add(c echo.Context) error {
	member, err := membership(h.member, c)
	if err != nil {
		return c.JSON(http.StatusForbidden, api.Error(err.Error()))
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect book"))
	}

	book, err := h.repository.Find(&id)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("book not found"))
	}

	m := model.DetermineBookObject(c.Request().URL.String())
	if err := c.Bind(m); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect book data"))
	}

	if err := c.Validate(m); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect book data", err.Error()))
	}

	m.SetCreatedBy(member.UserId)
	entity, err := h.repository.AddObject(book, m)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(entity, api.CreateRequest(c)))
}

func (h *Book) byOrganization(c echo.Context) error {
	member, err := membership(h.member, c)
	if err != nil {
		return c.JSON(http.StatusForbidden, api.Error(err.Error()))
	}

	request := api.CreateRequest(c)
	books, total, err := h.repository.FindManyByOrganization(member.OrganizationId, request)
	if err != nil {
		return c.JSON(http.StatusNotFound, api.Error("no books found in the organization"))
	}

	return c.JSON(http.StatusOK, api.Success(books, request, total))
}

// CreateBook creates a new instance of Book handler.
func CreateBook(g *echo.Group, d *sqlx.DB) *Book {
	return (&Book{echo: g, database: d}).initialize()
}
