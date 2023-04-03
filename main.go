package main

import (
	"flag"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lollipopkit/nano-db/api"
	. "github.com/lollipopkit/nano-db/cfg"
	"github.com/lollipopkit/nano-db/db"
	mid "github.com/lollipopkit/nano-db/middleware"
)

func main() {
	parseCli()

	initConfig()
	initCacher()
	initAcl()

	startWeb()
}

func parseCli() {
	userName := flag.String("u", "", "generate the cookie with -u <username>")
	dbName := flag.String("d", "", "update acl rules with -d <dbname>")
	flag.Parse()

	// generate cookie & update acl rules
	if *userName != "" {
		if *dbName == "" {
			println("[Cookie]\n ", api.GenCookie(*userName))
		} else {
			UpdateAcl(userName, dbName)
		}
		os.Exit(0)
	}
}

func startWeb() {
	e := echo.New()

	if Cfg.Log.Enable {
		e.Use(mid.Logger)
	}
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(Cfg.Security.RateLimit)))
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
	e.Start(Cfg.Addr)
}

func initCacher() {
	// Use these funcs to init Cacher
	// or params for cacher will be ignored
	// due to Golang `var init & func init()` sequence
	api.InitCacher()
	db.InitCacher()
}

func initAcl() {
	go func() {
		for {
			err := Acl.Load()
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Minute)
		}
	}()
}

func initConfig() {
	err := Cfg.Load()
	if err != nil {
		panic("AppConfig.Load(): " + err.Error())
	}
	if Cfg.Cache.ActiveRate < 0 || Cfg.Cache.ActiveRate > 1 {
		panic("invalid cache rate")
	}
}
