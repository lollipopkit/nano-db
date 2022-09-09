package middleware

import (
	"git.lolli.tech/lollipopkit/nano-db/consts"
	"github.com/labstack/echo/v4/middleware"
)

var Logger = middleware.LoggerWithConfig(middleware.LoggerConfig{
	Format:  consts.LogFormat,
	Skipper: consts.StaticLogSkipper,
})