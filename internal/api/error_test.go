package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

func TestCustomHTTPErrorHandler(t *testing.T) {
	e := echo.New()

	t.Run("StandardError", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := errors.New("some standard error")
		CustomHTTPErrorHandler(c, err)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var response Response
		_ = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, statusError, response.Status)
		assert.Equal(t, "internal server error", response.Message)
	})

	t.Run("HTTPError", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := echo.NewHTTPError(http.StatusNotFound, "resource not found")
		CustomHTTPErrorHandler(c, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

		var response Response
		_ = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, statusError, response.Status)
		assert.Equal(t, "resource not found", response.Message)
	})

	t.Run("HeadRequest", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodHead, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := echo.NewHTTPError(http.StatusBadRequest, "bad request")
		CustomHTTPErrorHandler(c, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Empty(t, rec.Body.String())
	})

	t.Run("ResponseAlreadyCommitted", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Commit the response
		c.Response().WriteHeader(http.StatusCreated)

		err := errors.New("some error after commit")
		CustomHTTPErrorHandler(c, err)

		// The status code should remain StatusCreated and not change to 500
		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}
