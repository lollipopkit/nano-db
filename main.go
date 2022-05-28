package main

import (
	"flag"

	"git.lolli.tech/LollipopKit/nano-db/api"
	"git.lolli.tech/LollipopKit/nano-db/consts"
	"git.lolli.tech/LollipopKit/nano-db/logger"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	addr := flag.String("u", "0.0.0.0:3777", "specific the addr to listen")
	userName := flag.String("c", "", "generate the cookie with -c <username>")
	flag.Parse()

	if consts.CookieSalt == "nano-db" {
		println(consts.CookieNotChanged)
	}

	// generate cookie
	if *userName != "" {
		println(api.GenCookie(*userName))
		return
	}

	// setup logger
	go logger.Setup()

	// Echo instance
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:  consts.LogFormat,
		Skipper: consts.StaticLogSkipper,
	}))
	e.Use(middleware.Recover())

	// Routes
	e.HEAD("/", api.Home)
	e.GET("/", api.Status)

	e.HEAD("/:db", api.Init)
	e.GET("/:db", api.Cols)
	e.DELETE("/:db", api.DeleteDB)

	e.GET("/:db/:col", api.IDs)
	e.DELETE("/:db/:col", api.DeleteCol)

	e.GET("/:db/:col/:id", api.Read)
	e.POST("/:db/:col/:id", api.Write)
	e.DELETE("/:db/:col/:id", api.Delete)

	// Start server
	e.HideBanner = true
	e.Logger.Fatal(e.Start(*addr))
}
