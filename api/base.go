package api

import (
	"github.com/labstack/echo/v4"
)

func permissionDenied(c echo.Context) error {
	return c.String(403, "permission denied")
}

// send auto add '.db/' to path
func send(c echo.Context, path string) error {
	return c.File(path)
}

func NotFound(c echo.Context) error {
	return c.NoContent(404)
}