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

// Book represents a book handler.
type Book struct {
	echo         *echo.Group
	database     *sqlx.DB
	repository   *database.BookRepository
	organization *database.OrganizationRepository
	member       *database.MemberRepository
}

func (h *Book) initialize() *Book {
	h.repository = database.CreateBook(h.database)
	h.organization = database.CreateOrganization(h.database)
	h.member = database.CreateMember(h.database)

	// h.echo.GET("/books/:id", h.one)
	// h.echo.GET("/books", h.many)
	h.echo.POST("/books", h.create)
	h.echo.POST("/books/:id/categories", h.add)
	h.echo.POST("/books/:id/locations", h.add)
	// h.echo.PUT("/books/:id", h.update)
	// h.echo.DELETE("/books/:id", h.delete)

	return h
}

func (h *Book) create(c echo.Context) error {
	book := &model.Book{}
	if err := c.Bind(book); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect book data"))
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(book); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect book data", err.Error()))
	}

	claims := token.FromContext(c)
	if claims.Id == nil {
		return c.JSON(http.StatusBadRequest, api.Error("invalid authentication token"))
	}

	member, err := h.member.Find(book.OrganizationId, claims.Id)
	if err != nil || member == nil {
		return c.JSON(http.StatusForbidden, api.Error("only organization members can create books"))
	}

	book.CreatedBy = claims.Id
	entity, err := h.repository.Create(book)
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

	m := model.DetermineBookObject(c.Request().URL.String())
	if err := c.Bind(m); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect book data"))
	}

	if err := c.Validate(m); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect book data", err.Error()))
	}

	claims := token.FromContext(c)
	if claims.Id == nil {
		return c.JSON(http.StatusBadRequest, api.Error(errorToken))
	}

	member, err := h.member.Find(book.OrganizationId, claims.Id)
	if err != nil || member == nil {
		return c.JSON(http.StatusForbidden, api.Error(errorMember))
	}

	m.SetCreatedBy(claims.Id)
	entity, err := h.repository.AddObject(book, m)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(entity, api.CreateRequest(c)))
}

// CreateBook creates a new instance of Book handler.
func CreateBook(g *echo.Group, d *sqlx.DB) *Book {
	return (&Book{echo: g, database: d}).initialize()
}
