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
	"github.com/jurgisjaska/binbogami/internal/database/location"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

type mockLocationRepository struct {
	locations      map[uuid.UUID]*location.Location
	failOnFind     bool
	failOnFindMany bool
}

func (m *mockLocationRepository) Find(id uuid.UUID) (*location.Location, error) {
	if m.failOnFind {
		return nil, errors.New("database error finding location")
	}
	l, ok := m.locations[id]
	if !ok || l.DeletedAt != nil {
		return nil, errors.New("location not found")
	}
	return l, nil
}

func (m *mockLocationRepository) FindMany(request *api.Request) (*location.Locations, int, error) {
	if m.failOnFindMany {
		return nil, 0, errors.New("database error finding locations")
	}
	var res location.Locations
	for _, l := range m.locations {
		if l.DeletedAt != nil {
			continue
		}
		res = append(res, *l)
	}
	return &res, len(res), nil
}

func (m *mockLocationRepository) ByBook(b *book.Book, id *uuid.UUID) (*location.Location, error) {
	return nil, nil
}

func (m *mockLocationRepository) ManyByBook(b *book.Book) (*location.Locations, error) {
	return nil, nil
}

func (m *mockLocationRepository) Create(l *models.Location) (*location.Location, error) {
	return nil, nil
}

func createLocationTestFixtures() map[uuid.UUID]*location.Location {
	jackID := uuid.MustParse("05e7257a-b21c-11ee-9a7a-5ab75f0c1cab")
	carterID := uuid.MustParse("1adcdaf6-b21c-11ee-9a7a-5ab75f0c1cab")

	now := time.Now()
	desc1 := "SGC Headquarters Base"
	desc2 := "Stargate Operation Room"
	desc3 := "Compromised by Anubis forces and abandoned"

	addr1 := "1 Norad Rd, Colorado Springs, CO 80906, USA"
	addr2 := "Sub-Level 28, Cheyenne Mountain Complex, CO, USA"
	addr3 := "P3X-984 Sector 2"

	loc1ID := uuid.MustParse("e0000000-0000-4000-8000-000000000001")
	loc2ID := uuid.MustParse("e0000000-0000-4000-8000-000000000002")
	locDeletedID := uuid.MustParse("e0000000-0000-4000-8000-000000000027")

	return map[uuid.UUID]*location.Location{
		loc1ID: {
			Id:          &loc1ID,
			Name:        "Cheyenne Mountain Complex",
			Description: &desc1,
			Address:     &addr1,
			CreatedBy:   &jackID,
			CreatedAt:   now,
		},
		loc2ID: {
			Id:          &loc2ID,
			Name:        "Gate Room (Sub-Level 28)",
			Description: &desc2,
			Address:     &addr2,
			CreatedBy:   &carterID,
			CreatedAt:   now,
		},
		locDeletedID: {
			Id:          &locDeletedID,
			Name:        "Decommissioned Alpha Site 1",
			Description: &desc3,
			Address:     &addr3,
			CreatedBy:   &jackID,
			CreatedAt:   now,
			DeletedAt:   &now,
		},
	}
}

func TestLocationIndex(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		targetURL      string
		mockRepo       *mockLocationRepository
		expectedStatus int
		expectInBody   []string
	}{
		{
			name:           "Default index returns active locations from fixtures",
			targetURL:      "/v1/locations",
			mockRepo:       &mockLocationRepository{locations: createLocationTestFixtures()},
			expectedStatus: http.StatusOK,
			expectInBody:   []string{"Cheyenne Mountain Complex", "Gate Room (Sub-Level 28)"},
		},
		{
			name:           "Database error on FindMany",
			targetURL:      "/v1/locations",
			mockRepo:       &mockLocationRepository{locations: createLocationTestFixtures(), failOnFindMany: true},
			expectedStatus: http.StatusNotFound,
			expectInBody:   []string{"database error finding locations"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			g := e.Group("/v1")
			h := &Location{
				echo:       g,
				repository: tt.mockRepo,
				auditlog:   logger,
			}
			g.GET("/locations", h.index)

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

func TestLocationShow(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name           string
		targetURL      string
		mockRepo       *mockLocationRepository
		expectedStatus int
		expectInBody   []string
	}{
		{
			name:           "Show existing active location fixture",
			targetURL:      "/v1/locations/e0000000-0000-4000-8000-000000000001",
			mockRepo:       &mockLocationRepository{locations: createLocationTestFixtures()},
			expectedStatus: http.StatusOK,
			expectInBody:   []string{"e0000000-0000-4000-8000-000000000001", "Cheyenne Mountain Complex"},
		},
		{
			name:           "Show non-existent location UUID",
			targetURL:      "/v1/locations/00000000-0000-0000-0000-000000000000",
			mockRepo:       &mockLocationRepository{locations: createLocationTestFixtures()},
			expectedStatus: http.StatusNotFound,
			expectInBody:   []string{"location not found"},
		},
		{
			name:           "Show deleted location fixture",
			targetURL:      "/v1/locations/e0000000-0000-4000-8000-000000000027",
			mockRepo:       &mockLocationRepository{locations: createLocationTestFixtures()},
			expectedStatus: http.StatusNotFound,
			expectInBody:   []string{"location not found"},
		},
		{
			name:           "Show with invalid UUID format",
			targetURL:      "/v1/locations/invalid-location-uuid",
			mockRepo:       &mockLocationRepository{locations: createLocationTestFixtures()},
			expectedStatus: http.StatusBadRequest,
			expectInBody:   []string{"incorrect location"},
		},
		{
			name:           "Database error on Find",
			targetURL:      "/v1/locations/e0000000-0000-4000-8000-000000000001",
			mockRepo:       &mockLocationRepository{locations: createLocationTestFixtures(), failOnFind: true},
			expectedStatus: http.StatusNotFound,
			expectInBody:   []string{"location not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			g := e.Group("/v1")
			h := &Location{
				echo:       g,
				repository: tt.mockRepo,
				auditlog:   logger,
			}
			g.GET("/locations/:id", h.show)

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
