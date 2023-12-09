package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lollipopkit/gommon/log"
	"github.com/lollipopkit/nano-db/cfg"
)

var (
	rateLimiterStore = middleware.NewRateLimiterMemoryStore(cfg.App.Security.RateLimit)
)

const (
	contextKeyPath     = "path"
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

func CheckPathAndPerm(depth uint8) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 检查路径
			var paths []string
			switch depth {
			case 1:
				paths = []string{c.Param("db")}
			case 2:
				paths = []string{c.Param("db"), c.Param("dir")}
			case 3:
				paths = []string{c.Param("db"), c.Param("dir"), c.Param("file")}
			default:
				return c.String(cePath, "invalid depth")
			}
			p, err := checkAndJoinPath(paths...)
			if err != nil {
				return c.String(cePath, err.Error())
			}
			c.Set(contextKeyPath, p)

			// 检查权限
			if !checkPermission(c) {
				return permissionDenied(c)
			}

			return next(c)
		}
	}
}
