package v1

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/model"
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
	member, err := membership(h.member, c)
	if err != nil {
		return c.JSON(http.StatusForbidden, api.Error(err.Error()))
	}

	entry := &model.Entry{}
	if err := c.Bind(entry); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect entry data", err.Error()))
	}

	if err := c.Validate(entry); err != nil {
		return c.JSON(http.StatusBadRequest, api.Errors("incorrect entry data", err.Error()))
	}

	book, err := h.book.Find(entry.BookId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error("invalid book for the entry"))
	}
	if book.ClosedAt != nil {
		return c.JSON(http.StatusBadRequest, api.Error("book closed"))
	}

	_, err = h.category.ByBook(book, entry.CategoryId)
	if err != nil {
		return c.JSON(http.StatusForbidden, api.Error("category does not belong to the book"))
	}

	_, err = h.location.ByBook(book, entry.LocationId)
	if err != nil {
		return c.JSON(http.StatusForbidden, api.Error("location does not belong to the book"))
	}

	entry.CreatedBy = member.UserId
	entity, err := h.repository.Create(entry)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, api.Success(entity, api.CreateRequest(c)))
}

func CreateEntry(g *echo.Group, d *sqlx.DB) *Entry {
	return (&Entry{echo: g, database: d}).initialize()
}
