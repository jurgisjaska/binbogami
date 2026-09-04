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
	"github.com/jurgisjaska/binbogami/internal/database/entry"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

type mockEntryRepository struct {
	entries        map[uuid.UUID]*entry.Entry
	failOnFind     bool
	failOnFindMany bool
}

func (m *mockEntryRepository) Find(id uuid.UUID) (*entry.Entry, error) {
	if m.failOnFind {
		return nil, errors.New("database error finding entry")
	}
	e, ok := m.entries[id]
	if !ok || e.DeletedAt != nil {
		return nil, errors.New("entry not found")
	}
	return e, nil
}

func (m *mockEntryRepository) FindMany(request *api.Request) (*entry.Entries, int, error) {
	if m.failOnFindMany {
		return nil, 0, errors.New("database error finding entries")
	}
	var res entry.Entries
	for _, e := range m.entries {
		if e.DeletedAt != nil {
			continue
		}
		res = append(res, *e)
	}
	return &res, len(res), nil
}

func (m *mockEntryRepository) Create(e *models.Entry) (*entry.Entry, error) {
	return nil, nil
}

func createEntryTestFixtures() map[uuid.UUID]*entry.Entry {
	carterID := uuid.MustParse("1adcdaf6-b21c-11ee-9a7a-5ab75f0c1cab")
	tealcID := uuid.MustParse("aff84550-b21f-11ee-8ac0-5ab75f0c1cab")

	bookID := uuid.MustParse("7b3a1f90-2c4d-4e5f-8a1b-9c0d1e2f3a4b")
	cat1ID := uuid.MustParse("f9e8d7c6-b5a4-4321-80f1-e2d3c4b5a697")
	cat2ID := uuid.MustParse("01234567-89ab-4cde-a012-34567890abcd")
	loc1ID := uuid.MustParse("e0000000-0000-4000-8000-000000000002")
	loc2ID := uuid.MustParse("e0000000-0000-4000-8000-000000000004")

	now := time.Now()
	desc1 := "P90 5.7x28mm ammunition crates (5000 rounds)"
	desc2 := "C4 plastic explosives block for off-world gate destruction"
	descDeleted := "SGC Infirmary emergency trauma surgical kits"

	entry1ID := uuid.MustParse("10000000-0000-4000-8000-000000000001")
	entry2ID := uuid.MustParse("10000000-0000-4000-8000-000000000002")
	entryDeletedID := uuid.MustParse("10000000-0000-4000-8000-000000000018")

	return map[uuid.UUID]*entry.Entry{
		entry1ID: {
			Id:          &entry1ID,
			Amount:      9600.42,
			Description: &desc1,
			BookId:      bookID,
			CategoryId:  &cat1ID,
			LocationId:  &loc1ID,
			CreatedBy:   &carterID,
			CreatedAt:   now,
		},
		entry2ID: {
			Id:          &entry2ID,
			Amount:      399.54,
			Description: &desc2,
			BookId:      bookID,
			CategoryId:  &cat2ID,
			LocationId:  &loc2ID,
			CreatedBy:   &tealcID,
			CreatedAt:   now,
		},
		entryDeletedID: {
			Id:          &entryDeletedID,
			Amount:      2353.31,
			Description: &descDeleted,
			BookId:      bookID,
			CategoryId:  &cat1ID,
			LocationId:  &loc1ID,
			CreatedBy:   &carterID,
			CreatedAt:   now,
			DeletedAt:   &now,
		},
	}
}

func TestEntryIndex(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		targetURL      string
		mockRepo       *mockEntryRepository
		expectedStatus int
		expectInBody   []string
	}{
		{
			name:           "Default index returns active entries from fixtures",
			targetURL:      "/v1/entries",
			mockRepo:       &mockEntryRepository{entries: createEntryTestFixtures()},
			expectedStatus: http.StatusOK,
			expectInBody:   []string{"P90 5.7x28mm ammunition crates", "C4 plastic explosives block"},
		},
		{
			name:           "Database error on FindMany",
			targetURL:      "/v1/entries",
			mockRepo:       &mockEntryRepository{entries: createEntryTestFixtures(), failOnFindMany: true},
			expectedStatus: http.StatusInternalServerError,
			expectInBody:   []string{"database error finding entries"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			g := e.Group("/v1")
			h := &Entry{
				echo:       g,
				repository: tt.mockRepo,
				auditlog:   logger,
			}
			g.GET("/entries", h.index)

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

func TestEntryShow(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		targetURL      string
		mockRepo       *mockEntryRepository
		expectedStatus int
		expectInBody   []string
	}{
		{
			name:           "Show existing active entry fixture",
			targetURL:      "/v1/entries/10000000-0000-4000-8000-000000000001",
			mockRepo:       &mockEntryRepository{entries: createEntryTestFixtures()},
			expectedStatus: http.StatusOK,
			expectInBody:   []string{"10000000-0000-4000-8000-000000000001", "P90 5.7x28mm ammunition crates"},
		},
		{
			name:           "Show non-existent entry UUID",
			targetURL:      "/v1/entries/00000000-0000-0000-0000-000000000000",
			mockRepo:       &mockEntryRepository{entries: createEntryTestFixtures()},
			expectedStatus: http.StatusNotFound,
			expectInBody:   []string{"entry not found"},
		},
		{
			name:           "Show deleted entry fixture",
			targetURL:      "/v1/entries/10000000-0000-4000-8000-000000000018",
			mockRepo:       &mockEntryRepository{entries: createEntryTestFixtures()},
			expectedStatus: http.StatusNotFound,
			expectInBody:   []string{"entry not found"},
		},
		{
			name:           "Show with invalid UUID format",
			targetURL:      "/v1/entries/invalid-entry-uuid",
			mockRepo:       &mockEntryRepository{entries: createEntryTestFixtures()},
			expectedStatus: http.StatusBadRequest,
			expectInBody:   []string{"incorrect entry"},
		},
		{
			name:           "Database error on Find",
			targetURL:      "/v1/entries/10000000-0000-4000-8000-000000000001",
			mockRepo:       &mockEntryRepository{entries: createEntryTestFixtures(), failOnFind: true},
			expectedStatus: http.StatusNotFound,
			expectInBody:   []string{"entry not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			g := e.Group("/v1")
			h := &Entry{
				echo:       g,
				repository: tt.mockRepo,
				auditlog:   logger,
			}
			g.GET("/entries/:id", h.show)

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
