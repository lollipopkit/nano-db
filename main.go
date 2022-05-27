package main

import (
	"flag"

	"github.com/LollipopKit/nano-db/api"
	"github.com/LollipopKit/nano-db/consts"
	"github.com/LollipopKit/nano-db/logger"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	addr := flag.String("u", "0.0.0.0:3777", "specific the addr to listen")
	userName := flag.String("c", "", "generate the cookie with -c <username>")
	flag.Parse()

	// setup logger
	go logger.Setup()

	// generate cookie
	if *userName != "" {
		println(api.GenCookie(*userName))
		return
	}

	// Echo instance
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:  consts.LogFormat,
		Skipper: consts.StaticLogSkipper,
	}))
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", api.Home)
	e.GET("/status", api.Status)
	e.GET("/:db/:col/:id", api.Read)
	e.POST("/:db/:col/:id", api.Write)
	e.DELETE("/:db/:col/:id", api.Delete)

	// Start server
	e.HideBanner = true
	e.Logger.Fatal(e.Start(*addr))
}
