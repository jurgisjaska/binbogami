package api

import (
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	defaultPage     int    = 1
	defaultLimit    int    = 25
	defaultOrder    string = "asc"
	orderAscending  string = "asc"
	orderDescending string = "desc"
	statusError     string = "error"
	statusSuccess   string = "success"
)

type (
	// Response represents the structure for the response.
	Response struct {
		Status   string           `json:"status"`
		Data     interface{}      `json:"data,omitempty"`
		Message  string           `json:"message,omitempty"`
		Metadata ResponseMetadata `json:"metadata"`
	}

	// ResponseMetadata represents the metadata for the response.
	ResponseMetadata struct {
		Total int     `json:"total"`
		Limit int     `json:"limit"`
		Page  int     `json:"page"`
		Pages float64 `json:"pages"`
	}

	// Request represents the structure for making a request.
	Request struct {
		Page    int    `json:"page"`
		Limit   int    `json:"limit"`
		OrderBy string `json:"order_by"`
		Order   string `json:"order"`
	}
)

// CreateRequest creates a new Request object based on the query parameters in the provided echo.Context.
func CreateRequest(c echo.Context) *Request {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	orderBy := strings.ToLower(c.QueryParam("order_by"))
	order := strings.ToLower(c.QueryParam("order"))

	if page < 1 {
		page = defaultPage
	}

	if limit < 1 {
		limit = defaultLimit
	}

	if order != orderAscending && order != orderDescending {
		order = defaultOrder
	}

	return &Request{page, limit, orderBy, order}
}

// @todo all list pages MUST be updated with correct total!
func Success(data interface{}, req *Request, t ...int) *Response {
	total, pages := totalPages(data, req, t)

	return &Response{
		Status: statusSuccess,
		Data:   data,
		Metadata: ResponseMetadata{
			Total: total,
			Limit: req.Limit,
			Page:  req.Page,
			Pages: pages,
		},
	}
}

func totalPages(data interface{}, req *Request, t []int) (int, float64) {
	total := 1
	if len(t) > 0 && t[0] > 0 {
		// total is calculated during database operations
		total = t[0]
	} else {
		// total is calculated by response data
		dv := reflect.ValueOf(data)
		if dv.Kind() == reflect.Pointer {
			dv = dv.Elem()
		}

		if dv.Kind() == reflect.Slice || dv.Kind() == reflect.Array {
			total = dv.Len()
		}
	}

	pages := math.Ceil(float64(total) / float64(req.Limit))
	if pages == 0 {
		pages = 1
	}

	return total, pages
}

func Error(m string) *Response {
	return &Response{
		Status:  statusError,
		Message: strings.ToLower(m),
	}
}

func Errors(m string, e interface{}) *Response {
	// @todo better validation error responses
	// verify that errors are from validation
	// range trough them and build better response

	return &Response{
		Status:  statusError,
		Message: strings.ToLower(m),
		Data:    e,
	}
}
