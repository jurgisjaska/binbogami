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
	h.echo.PUT("/books", h.update)
	h.echo.GET("/books", h.byOrganization)
	h.echo.GET("/books/:id", h.show)
	h.echo.POST("/books/:id/categories", h.add)
	h.echo.POST("/books/:id/locations", h.add)

	return h
}

func (h *Book) create(c echo.Context) error {
	member, err := membership(h.member, c)
	if err != nil {
		return c.JSON(http.StatusForbidden, api.Error(err.Error()))
	}

	bm := &model.CreateBook{}
	bm.CreatedBy = member.UserId
	bm.OrganizationId = member.OrganizationId

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
	_, err := membership(h.member, c)
	if err != nil {
		return c.JSON(http.StatusForbidden, api.Error(err.Error()))
	}

	bm := &model.UpdateBook{}
	if err := c.Bind(bm); err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("incorrect book data"))
	}

	// user performing the change must be a member of the new organization
	// @todo this may not work if in the future there is an administrator dashboard
	_, err = verifyMembership(h.member, c, bm.OrganizationId)
	if err != nil {
		return c.JSON(http.StatusForbidden, api.Error("not a member of new organization"))
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

// CreateBook creates a new instance of Book handler.
func CreateBook(g *echo.Group, d *sqlx.DB) *Book {
	return (&Book{echo: g, database: d}).initialize()
}
