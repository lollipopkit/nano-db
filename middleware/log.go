package middleware

import (
	"github.com/labstack/echo/v4/middleware"
	"github.com/lollipopkit/nano-db/consts"
)

var Logger = middleware.LoggerWithConfig(middleware.LoggerConfig{
	Format:  consts.LogFormat,
	Skipper: consts.StaticLogSkipper,
})
