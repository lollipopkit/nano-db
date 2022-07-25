package consts

import (
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	LogFormat = "${time_rfc3339} ${method} ${status} ${uri} \nLatency: ${latency_human}  ${error}\n"
)

func StaticLogSkipper(context echo.Context) bool {
	url := context.Request().URL.Path
	if strings.Contains(url, "/static") || url == "/favicon.ico" {
		return true
	}
	return false
}
