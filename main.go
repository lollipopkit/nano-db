package main

import (
	"flag"

	. "git.lolli.tech/lollipopkit/nano-db/acl"
	"git.lolli.tech/lollipopkit/nano-db/api"
	"git.lolli.tech/lollipopkit/nano-db/consts"
	mid "git.lolli.tech/lollipopkit/nano-db/middleware"
	"git.lolli.tech/lollipopkit/nano-db/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	addr := flag.String("a", "0.0.0.0:3777", "specific the addr to listen")
	userName := flag.String("u", "", "generate the cookie with -u <username>")
	dbName := flag.String("d", "", "update acl rules with -d <dbname>")
	cacheLen := flag.Uint("l", 100, "set the max length of cache")
	log := flag.Bool("log", false, "enable log")
	flag.Parse()

	utils.InitSalt()

	consts.CacherMaxLength = *cacheLen

	// generate cookie & update acl rules
	if *userName != "" {
		if *dbName == "" {
			println("[Cookie]\n ", api.GenCookie(*userName))
		} else {
			UpdateAcl(userName, dbName)
		}
		return
	}

	startHttp(addr, *log)
}

func startHttp(addr *string, log bool) {
	// Echo instance
	e := echo.New()

	if log {
		e.Use(mid.Logger)
	}
	e.Use(middleware.Recover())

	// Routes
	e.HEAD("/", api.Alive)
	e.GET("/", api.Status)

	e.GET("/:db", api.Dirs)
	e.DELETE("/:db", api.DeleteDB)
	e.POST("/:db", api.SearchInDB)

	e.GET("/:db/:dir", api.Files)
	e.DELETE("/:db/:dir", api.DeleteCol)
	e.POST("/:db/:dir", api.SearchInDir)

	e.GET("/:db/:dir/:file", api.Read)
	e.POST("/:db/:dir/:file", api.Write)
	e.DELETE("/:db/:dir/:file", api.Delete)

	// Start server
	e.HideBanner = true
	e.Start(*addr)
}
