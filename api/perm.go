package api

import (
	"github.com/labstack/echo/v4"
	"github.com/lollipopkit/nano-db/cfg"
	"github.com/lollipopkit/nano-db/cst"
)

func checkPermission(c echo.Context, action, dbName string) bool {
	sn := c.Request().Header.Get(cst.HeaderKey)
	if len(sn) == 0 {
		return false
	}

	return cfg.Acl.Can(dbName, sn)
}
