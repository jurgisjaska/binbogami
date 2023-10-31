package api

import (
	"math"
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
	Response struct {
		Status   string           `json:"status"`
		Data     interface{}      `json:"data,omitempty"`
		Message  string           `json:"message,omitempty"`
		Metadata ResponseMetadata `json:"metadata"`
	}

	ResponseMetadata struct {
		Total int     `json:"total"`
		Limit int     `json:"limit"`
		Page  int     `json:"page"`
		Pages float64 `json:"pages"`
	}

	Request struct {
		Page    int    `json:"page"`
		Limit   int    `json:"limit"`
		OrderBy string `json:"order_by"`
		Order   string `json:"order"`
	}
)

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

func Success(d interface{}, t int, r *Request) *Response {
	pages := math.Ceil(float64(t / r.Limit))
	if pages == 0 {
		pages = 1
	}

	return &Response{
		Status: statusSuccess,
		Data:   d,
		Metadata: ResponseMetadata{
			Total: t,
			Limit: r.Limit,
			Page:  r.Page,
			Pages: pages,
		},
	}
}

func Error(m string) *Response {
	return &Response{
		Status:  statusError,
		Message: strings.ToLower(m),
	}
}

func Errors(m string, e interface{}) *Response {
	return &Response{
		Status:  statusError,
		Message: strings.ToLower(m),
		Data:    e,
	}
}
