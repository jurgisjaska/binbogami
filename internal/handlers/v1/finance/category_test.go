package finance

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jurgisjaska/binbogami/internal/api"
	"github.com/jurgisjaska/binbogami/internal/api/models"
	"github.com/jurgisjaska/binbogami/internal/database/book"
	"github.com/jurgisjaska/binbogami/internal/database/category"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

type mockCategoryRepository struct {
	categories     map[uuid.UUID]*category.Category
	failOnFind     bool
	failOnFindMany bool
}

func (m *mockCategoryRepository) Find(id uuid.UUID) (*category.Category, error) {
	if m.failOnFind {
		return nil, errors.New("database error finding category")
	}
	c, ok := m.categories[id]
	if !ok || c.DeletedAt != nil {
		return nil, errors.New("category not found")
	}
	return c, nil
}

func (m *mockCategoryRepository) FindMany(request *api.Request) (*category.Categories, int, error) {
	if m.failOnFindMany {
		return nil, 0, errors.New("database error finding categories")
	}
	var res category.Categories
	for _, c := range m.categories {
		if c.DeletedAt != nil {
			continue
		}
		res = append(res, *c)
	}
	return &res, len(res), nil
}

func (m *mockCategoryRepository) ByBook(b *book.Book, id *uuid.UUID) (*category.Category, error) {
	return nil, nil
}

func (m *mockCategoryRepository) ManyByBook(b *book.Book) (*category.Categories, error) {
	return nil, nil
}

func (m *mockCategoryRepository) Create(c *models.Category) (*category.Category, error) {
	return nil, nil
}

func (m *mockCategoryRepository) Remove(c *category.Category) error {
	return nil
}

func createCategoryTestFixtures() map[uuid.UUID]*category.Category {
	jackID := uuid.MustParse("05e7257a-b21c-11ee-9a7a-5ab75f0c1cab")
	danielID := uuid.MustParse("2b63b228-b21c-11ee-9a7a-5ab75f0c1cab")

	now := time.Now()
	desc1 := "P90 ammo, C4, and standard issue gear"
	desc2 := "Tools and resources for archaeological analysis"
	desc3 := "Off-book acquisitions (Unauthorized)"

	color1 := "#4CAF50"
	color2 := "#2196F3"
	color3 := "#000000"

	cat1ID := uuid.MustParse("d4c3b2a1-0e9f-48d7-b6c5-a4b3c2d1e0f9")
	cat2ID := uuid.MustParse("f9e8d7c6-b5a4-4321-80f1-e2d3c4b5a697")
	catDeletedID := uuid.MustParse("fedcba98-7654-4321-8fed-cba987654321")

	return map[uuid.UUID]*category.Category{
		cat1ID: {
			Id:          &cat1ID,
			Name:        "Mission Supplies",
			Description: &desc1,
			Color:       &color1,
			CreatedBy:   &jackID,
			CreatedAt:   now,
		},
		cat2ID: {
			Id:          &cat2ID,
			Name:        "Artifact Research",
			Description: &desc2,
			Color:       &color2,
			CreatedBy:   &danielID,
			CreatedAt:   now,
		},
		catDeletedID: {
			Id:          &catDeletedID,
			Name:        "NID Black Budget",
			Description: &desc3,
			Color:       &color3,
			CreatedBy:   &jackID,
			CreatedAt:   now,
			DeletedAt:   &now,
		},
	}
}

func TestCategoryIndex(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		targetURL      string
		mockRepo       *mockCategoryRepository
		expectedStatus int
		expectInBody   []string
	}{
		{
			name:           "Default index returns active categories from fixtures",
			targetURL:      "/v1/categories",
			mockRepo:       &mockCategoryRepository{categories: createCategoryTestFixtures()},
			expectedStatus: http.StatusOK,
			expectInBody:   []string{"Mission Supplies", "Artifact Research"},
		},
		{
			name:           "Database error on FindMany",
			targetURL:      "/v1/categories",
			mockRepo:       &mockCategoryRepository{categories: createCategoryTestFixtures(), failOnFindMany: true},
			expectedStatus: http.StatusNotFound,
			expectInBody:   []string{"database error finding categories"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			g := e.Group("/v1")
			h := &Category{
				echo:       g,
				repository: tt.mockRepo,
				auditlog:   logger,
			}
			g.GET("/categories", h.index)

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

func TestCategoryShow(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		targetURL      string
		mockRepo       *mockCategoryRepository
		expectedStatus int
		expectInBody   []string
	}{
		{
			name:           "Show existing active category fixture",
			targetURL:      "/v1/categories/d4c3b2a1-0e9f-48d7-b6c5-a4b3c2d1e0f9",
			mockRepo:       &mockCategoryRepository{categories: createCategoryTestFixtures()},
			expectedStatus: http.StatusOK,
			expectInBody:   []string{"d4c3b2a1-0e9f-48d7-b6c5-a4b3c2d1e0f9", "Mission Supplies"},
		},
		{
			name:           "Show non-existent category UUID",
			targetURL:      "/v1/categories/00000000-0000-0000-0000-000000000000",
			mockRepo:       &mockCategoryRepository{categories: createCategoryTestFixtures()},
			expectedStatus: http.StatusNotFound,
			expectInBody:   []string{"category not found"},
		},
		{
			name:           "Show deleted category fixture",
			targetURL:      "/v1/categories/fedcba98-7654-4321-8fed-cba987654321",
			mockRepo:       &mockCategoryRepository{categories: createCategoryTestFixtures()},
			expectedStatus: http.StatusNotFound,
			expectInBody:   []string{"category not found"},
		},
		{
			name:           "Show with invalid UUID format",
			targetURL:      "/v1/categories/not-a-uuid",
			mockRepo:       &mockCategoryRepository{categories: createCategoryTestFixtures()},
			expectedStatus: http.StatusBadRequest,
			expectInBody:   []string{"incorrect category"},
		},
		{
			name:           "Database error on Find",
			targetURL:      "/v1/categories/d4c3b2a1-0e9f-48d7-b6c5-a4b3c2d1e0f9",
			mockRepo:       &mockCategoryRepository{categories: createCategoryTestFixtures(), failOnFind: true},
			expectedStatus: http.StatusNotFound,
			expectInBody:   []string{"category not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			g := e.Group("/v1")
			h := &Category{
				echo:       g,
				repository: tt.mockRepo,
				auditlog:   logger,
			}
			g.GET("/categories/:id", h.show)

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
