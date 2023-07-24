package api

import (
	"github.com/labstack/echo/v4"
	"github.com/lollipopkit/gommon/log"
)

func permissionDenied(c echo.Context) error {
	return c.String(403, "permission denied")
}

// send auto add '.db/' to path
func send(c echo.Context, path string) error {
	return c.File(path)
}

// Copy from echo.DefaultHTTPErrorHandler
func HandleErr(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	he, ok := err.(*echo.HTTPError)
	if ok {
		if he.Internal != nil {
			if herr, ok := he.Internal.(*echo.HTTPError); ok {
				he = herr
			}
		}
	} else {
		he = &echo.HTTPError{
			Code:    500,
			Message: "internal server error",
		}
	}

	err = c.NoContent(he.Code)
	if err != nil {
		log.Err(err.Error())
	}
}
