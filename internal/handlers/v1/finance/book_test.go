package finance

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/models"
	"github.com/jurgisjaska/binbogami/internal/database/book"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

type mockBookRepository struct {
	books            map[uuid.UUID]*book.Book
	failOnFind       bool
	failOnFindMany   bool
	failOnFindByName bool
}

func (m *mockBookRepository) Find(id uuid.UUID) (*book.Book, error) {
	if m.failOnFind {
		return nil, errors.New("database error finding book")
	}
	b, ok := m.books[id]
	if !ok || b.DeletedAt != nil {
		return nil, errors.New("book not found")
	}
	return b, nil
}

func (m *mockBookRepository) FindMany(request *api.Request, status string) (*book.Books, int, error) {
	if m.failOnFindMany {
		return nil, 0, errors.New("database error finding many books")
	}
	var res book.Books
	for _, b := range m.books {
		if b.DeletedAt != nil {
			continue
		}
		if status == "closed" && b.ClosedAt == nil {
			continue
		}
		if (status == "" || status == "active") && b.ClosedAt != nil {
			continue
		}
		res = append(res, *b)
	}
	return &res, len(res), nil
}

func (m *mockBookRepository) FindManyByName(req *api.Request, status string, search string) (*book.Books, int, error) {
	if m.failOnFindByName {
		return nil, 0, errors.New("database error finding books by name")
	}
	var res book.Books
	for _, b := range m.books {
		if b.DeletedAt != nil {
			continue
		}
		if strings.Contains(strings.ToLower(b.Name), strings.ToLower(search)) {
			res = append(res, *b)
		}
	}
	return &res, len(res), nil
}

func (m *mockBookRepository) Create(b *book.Book) error {
	return nil
}

func (m *mockBookRepository) Update(b *book.Book) error {
	return nil
}

func (m *mockBookRepository) AddObject(b *book.Book, model models.BookObject) (any, error) {
	return nil, nil
}

func createBookTestFixtures() map[uuid.UUID]*book.Book {
	jackID := uuid.MustParse("05e7257a-b21c-11ee-9a7a-5ab75f0c1cab")
	carterID := uuid.MustParse("1adcdaf6-b21c-11ee-9a7a-5ab75f0c1cab")
	tealcID := uuid.MustParse("aff84550-b21f-11ee-8ac0-5ab75f0c1cab")

	now := time.Now()
	desc2025 := "The book for year of 2025"
	desc2026 := "The book for year of 2026"
	descOffWorld := "Financial ledger for off-world trading, naquadah exchanges, and planetary bartering"
	descDeleted := "Book that was created incorrect and deleted"

	book2026ID := uuid.MustParse("7b3a1f90-2c4d-4e5f-8a1b-9c0d1e2f3a4b")
	book2025ID := uuid.MustParse("5e6f7a8b-9c0d-4e1f-b2a3-4b5c6d7e8f9a")
	bookOffWorldID := uuid.MustParse("9c8d7e6f-5a4b-4c2d-9e0f-9a8b7c6d5e4f")
	bookDeletedID := uuid.MustParse("1a2b3c4d-5e6f-4789-a0b1-c2d3e4f5a6b7")

	return map[uuid.UUID]*book.Book{
		book2026ID: {
			Id:          book2026ID,
			Name:        "Year 2026",
			Description: &desc2026,
			CreatedBy:   jackID,
			CreatedAt:   now,
		},
		book2025ID: {
			Id:          book2025ID,
			Name:        "Year 2025",
			Description: &desc2025,
			CreatedBy:   jackID,
			CreatedAt:   now,
			ClosedAt:    &now,
		},
		bookOffWorldID: {
			Id:          bookOffWorldID,
			Name:        "Year 2026 - Off-World Trade & Barter",
			Description: &descOffWorld,
			CreatedBy:   tealcID,
			CreatedAt:   now,
		},
		bookDeletedID: {
			Id:          bookDeletedID,
			Name:        "Year 2025 deleted",
			Description: &descDeleted,
			CreatedBy:   carterID,
			CreatedAt:   now,
			DeletedAt:   &now,
		},
	}
}

func TestBookIndex(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		targetURL      string
		mockRepo       *mockBookRepository
		expectedStatus int
		expectInBody   []string
	}{
		{
			name:           "Default index returns active books from fixtures",
			targetURL:      "/v1/books",
			mockRepo:       &mockBookRepository{books: createBookTestFixtures()},
			expectedStatus: http.StatusOK,
			expectInBody:   []string{"Year 2026", "Off-World"},
		},
		{
			name:           "Filter index by status closed",
			targetURL:      "/v1/books?status=closed",
			mockRepo:       &mockBookRepository{books: createBookTestFixtures()},
			expectedStatus: http.StatusOK,
			expectInBody:   []string{"Year 2025"},
		},
		{
			name:           "Filter index by status any",
			targetURL:      "/v1/books?status=any",
			mockRepo:       &mockBookRepository{books: createBookTestFixtures()},
			expectedStatus: http.StatusOK,
			expectInBody:   []string{"Year 2026", "Year 2025", "Off-World"},
		},
		{
			name:           "Search index by query parameter",
			targetURL:      "/v1/books?query=Off-World",
			mockRepo:       &mockBookRepository{books: createBookTestFixtures()},
			expectedStatus: http.StatusOK,
			expectInBody:   []string{"Off-World"},
		},
		{
			name:           "Database error on FindMany",
			targetURL:      "/v1/books",
			mockRepo:       &mockBookRepository{books: createBookTestFixtures(), failOnFindMany: true},
			expectedStatus: http.StatusNotFound,
			expectInBody:   []string{"no books found"},
		},
		{
			name:           "Database error on FindManyByName",
			targetURL:      "/v1/books?query=Search",
			mockRepo:       &mockBookRepository{books: createBookTestFixtures(), failOnFindByName: true},
			expectedStatus: http.StatusNotFound,
			expectInBody:   []string{"no books found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			g := e.Group("/v1")
			h := &Book{
				echo:       g,
				repository: tt.mockRepo,
				auditlog:   logger,
			}
			g.GET("/books", h.index)

			req := httptest.NewRequest(http.MethodGet, tt.targetURL, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			for _, exp := range tt.expectInBody {
				assert.Contains(t, rec.Body.String(), exp)
			}
		})
	}
}

func TestBookShow(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		targetURL      string
		mockRepo       *mockBookRepository
		expectedStatus int
		expectInBody   []string
	}{
		{
			name:           "Show existing active book fixture",
			targetURL:      "/v1/books/7b3a1f90-2c4d-4e5f-8a1b-9c0d1e2f3a4b",
			mockRepo:       &mockBookRepository{books: createBookTestFixtures()},
			expectedStatus: http.StatusOK,
			expectInBody:   []string{"7b3a1f90-2c4d-4e5f-8a1b-9c0d1e2f3a4b", "Year 2026"},
		},
		{
			name:           "Show non-existent book UUID",
			targetURL:      "/v1/books/00000000-0000-0000-0000-000000000000",
			mockRepo:       &mockBookRepository{books: createBookTestFixtures()},
			expectedStatus: http.StatusNotFound,
			expectInBody:   []string{"no books found"},
		},
		{
			name:           "Show deleted book fixture",
			targetURL:      "/v1/books/1a2b3c4d-5e6f-4789-a0b1-c2d3e4f5a6b7",
			mockRepo:       &mockBookRepository{books: createBookTestFixtures()},
			expectedStatus: http.StatusNotFound,
			expectInBody:   []string{"no books found"},
		},
		{
			name:           "Show with invalid UUID format",
			targetURL:      "/v1/books/invalid-uuid-format",
			mockRepo:       &mockBookRepository{books: createBookTestFixtures()},
			expectedStatus: http.StatusBadRequest,
			expectInBody:   []string{"incorrect book"},
		},
		{
			name:           "Database error on Find",
			targetURL:      "/v1/books/7b3a1f90-2c4d-4e5f-8a1b-9c0d1e2f3a4b",
			mockRepo:       &mockBookRepository{books: createBookTestFixtures(), failOnFind: true},
			expectedStatus: http.StatusNotFound,
			expectInBody:   []string{"no books found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			g := e.Group("/v1")
			h := &Book{
				echo:       g,
				repository: tt.mockRepo,
				auditlog:   logger,
			}
			g.GET("/books/:id", h.show)

			req := httptest.NewRequest(http.MethodGet, tt.targetURL, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			for _, exp := range tt.expectInBody {
				assert.Contains(t, rec.Body.String(), exp)
			}
		})
	}
}
