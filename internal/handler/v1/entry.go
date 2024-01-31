package v1

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/model"
	"github.com/jurgisjaska/binbogami/internal/api/token"
	"github.com/jurgisjaska/binbogami/internal/database"
	"github.com/labstack/echo/v4"
)

type Entry struct {
	echo       *echo.Group
	database   *sqlx.DB
	repository *database.EntryRepository

	member   *database.MemberRepository
	book     *database.BookRepository
	category *database.CategoryRepository
	location *database.LocationRepository
}

func (h *Entry) initialize() *Entry {
	h.repository = database.CreateEntry(h.database)
	h.member = database.CreateMember(h.database)
	h.book = database.CreateBook(h.database)
	h.category = database.CreateCategory(h.database)
	h.location = database.CreateLocation(h.database)

	h.echo.POST("/entries", h.create)

	return h
}

func (h *Entry) create(c echo.Context) error {
	entry := &model.Entry{}
	if err := c.Bind(entry); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect entry data", err.Error()))
	}

	if err := c.Validate(entry); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect entry data", err.Error()))
	}

	claims := token.FromContext(c)
	if claims.Id == nil {
		return c.JSON(http.StatusBadRequest, api.Error("invalid authentication token"))
	}

	book, err := h.book.Find(entry.BookId)
	if err != nil || book == nil {
		return c.JSON(http.StatusBadRequest, api.Error("invalid book for the entry"))
	}
	if book.ClosedAt != nil {
		return c.JSON(http.StatusBadRequest, api.Error("book closed"))
	}

	member, err := h.member.ByBook(book, claims.Id)
	if err != nil || member == nil {
		return c.JSON(http.StatusForbidden, api.Error("only organization members can create entries"))
	}

	category, err := h.category.ByBook(book, entry.CategoryId)
	if err != nil || category == nil {
		return c.JSON(http.StatusForbidden, api.Error("category does not belong to the book"))
	}

	location, err := h.location.ByBook(book, entry.LocationId)
	if err != nil || location == nil {
		return c.JSON(http.StatusForbidden, api.Error("location does not belong to the book"))
	}

	entry.CreatedBy = claims.Id
	entity, err := h.repository.Create(entry)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(entity, api.CreateRequest(c)))

	return nil
}

func CreateEntry(g *echo.Group, d *sqlx.DB) *Entry {
	return (&Entry{echo: g, database: d}).initialize()
}
