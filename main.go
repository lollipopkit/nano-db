package main

import (
	"flag"

	"git.lolli.tech/lollipopkit/nano-db/api"
	"git.lolli.tech/lollipopkit/nano-db/consts"
	"git.lolli.tech/lollipopkit/nano-db/logger"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	addr := flag.String("a", "0.0.0.0:3777", "specific the addr to listen")
	userName := flag.String("u", "", "generate the cookie with -n <username>")
	dbName := flag.String("d", "", "update acl rules with -d <dbname>")
	salt := flag.String("s", "", "set salt for cookie")
	cacheLen := flag.Int("l", 100, "set the max length of cache")
	flag.Parse()

	if *salt != "" {
		consts.CookieSalt = *salt
	}

	if consts.CookieSalt == "nano-db" {
		println(consts.CookieNotChanged)
	}

	consts.CacherMaxLength = *cacheLen

	// generate cookie & update acl rules
	if *userName != "" {
		if *dbName == "" {
			println("[Cookie]\n ", api.GenCookie(*userName))
		} else {
			updateAcl(userName, dbName)
		}
		return
	}

	// setup logger
	go logger.Setup()

	startWeb(addr)
}

func startWeb(addr *string) {
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

	e.GET("/:db", api.Cols)
	e.DELETE("/:db", api.DeleteDB)

	e.GET("/:db/:col", api.IDs)
	e.DELETE("/:db/:col", api.DeleteCol)

	e.HEAD("/:db/:col/:id", api.Exist)
	e.GET("/:db/:col/:id", api.Read)
	e.POST("/:db/:col/:id", api.Write)
	e.DELETE("/:db/:col/:id", api.Delete)

	// Start server
	e.HideBanner = true
	e.Logger.Fatal(e.Start(*addr))
}

func updateAcl(userName, dbName *string) {
	print("[ACL]\n  ")
	api.AclLock.RLock()
	if api.Acl.HaveDB(*dbName) {
		if !api.Acl.Can(*dbName, *userName) {
			api.AclLock.RUnlock()
			println("this db already initialized by other user")
			return
		}
		api.AclLock.RUnlock()
		println("you already initialized this db")
		return
	}

	api.AclLock.RUnlock()
	api.AclLock.Lock()
	err := api.Acl.UpdateRule(*dbName, *userName)
	api.AclLock.Unlock()

	if err != nil {
		println("[api.Init] acl.UpdateRule(): " + err.Error())
	} else {
		println("[api.Init] acl.UpdateRule(): success")
	}
}
