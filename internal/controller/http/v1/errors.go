package v1

import (
	"errors"
	"github.com/labstack/echo/v4"
)

var (
	ErrInvalidAuthHeader = errors.New("invalid auth header")
	ErrCannotParseToken  = errors.New("cannot parse token")
)

func newErrorResponse(c echo.Context, code int, message string) {
	_ = c.JSON(code, map[string]string{"errors": message})
}
