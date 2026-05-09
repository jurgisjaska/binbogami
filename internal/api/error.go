package api

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
)

func CustomHTTPErrorHandler(c *echo.Context, err error) {
	if resp, uErr := echo.UnwrapResponse(c.Response()); uErr == nil {
		if resp.Committed {
			return
		}
	}

	code := http.StatusInternalServerError
	var message = http.StatusText(http.StatusInternalServerError)

	var he *echo.HTTPError
	if errors.As(err, &he) {
		code = he.Code
		message = he.Message
	}

	if c.Request().Method == http.MethodHead {
		_ = c.NoContent(code)
		return
	}

	c.Response().Header().Set("Content-Type", "application/json")
	_ = c.JSON(code, Error(message))
}
