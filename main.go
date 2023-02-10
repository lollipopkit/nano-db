package main

import (
	"flag"

	. "git.lolli.tech/lollipopkit/nano-db/acl"
	"git.lolli.tech/lollipopkit/nano-db/api"
	"git.lolli.tech/lollipopkit/nano-db/consts"
	"git.lolli.tech/lollipopkit/nano-db/db"
	mid "git.lolli.tech/lollipopkit/nano-db/middleware"
	"git.lolli.tech/lollipopkit/nano-db/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	addr := flag.String("a", "0.0.0.0:3777", "specific the addr to listen")
	userName := flag.String("u", "", "generate the cookie with -u <username>")
	dbName := flag.String("d", "", "update acl rules with -d <dbname>")
	cacheLen := flag.Int("l", 100, "set the max length of cache")
	cacheRate := flag.Float64("r", 0.8, "set the activeRate of cacher (0.0-1.0)")
	log := flag.Bool("log", false, "enable log")
	flag.Parse()

	utils.InitSalt()

	consts.CacherMaxLength = *cacheLen
	if *cacheRate < 0 || *cacheRate > 1 {
		println("invalid cache rate")
	}
	consts.CacherActiveRate = *cacheRate
	// Use these funcs to init Cacher
	// or params for cacher will be ignored
	// due to Golang `var init & init func` sequence
	api.InitCacher()
	db.InitCacher()

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
	e.POST("/:db", api.SearchDB)

	e.GET("/:db/:dir", api.Files)
	e.DELETE("/:db/:dir", api.DeleteDir)
	e.POST("/:db/:dir", api.SearchDir)

	e.GET("/:db/:dir/:file", api.Read)
	e.POST("/:db/:dir/:file", api.Write)
	e.DELETE("/:db/:dir/:file", api.Delete)

	// Start server
	e.HideBanner = true
	e.Start(*addr)
}
