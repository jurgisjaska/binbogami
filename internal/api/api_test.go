package api

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected *Response
	}{
		{
			name:     "empty input",
			input:    "",
			expected: &Response{Status: statusError, Message: ""},
		},
		{
			name:     "normal case",
			input:    "ERROR MESSAGE",
			expected: &Response{Status: statusError, Message: "error message"},
		},
		{
			name:     "with numbers",
			input:    "ERROR123 MESSAGE",
			expected: &Response{Status: statusError, Message: "error123 message"},
		},
		{
			name:     "with special characters",
			input:    "ERROR!@# MESSAGE",
			expected: &Response{Status: statusError, Message: "error!@# message"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := Error(tc.input)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Expected: %v, Actual: %v", tc.expected, actual)
			}
		})
	}
}

func TestCreateRequest(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name       string
		urlParams  string
		wantResult *Request
	}{
		{
			name:       "DefaultValues",
			urlParams:  "?page=&limit=&order_by=&order=",
			wantResult: &Request{defaultPage, defaultLimit, "", defaultOrder},
		},
		{
			name:       "PageAndLimit",
			urlParams:  "?page=2&limit=5",
			wantResult: &Request{2, 5, "", defaultOrder},
		},
		{
			name:       "OrderByAndOrder",
			urlParams:  "?order_by=name&order=asc",
			wantResult: &Request{defaultPage, defaultLimit, "name", "asc"},
		},
		{
			name:       "InvalidOrder",
			urlParams:  "?order=invalid",
			wantResult: &Request{defaultPage, defaultLimit, "", defaultOrder},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/"+tt.urlParams, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			gotResult := CreateRequest(c)

			assert.Equal(t, tt.wantResult, gotResult)
		})
	}
}
