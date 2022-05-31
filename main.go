package main

import (
	"flag"
	"io/ioutil"
	"math/rand"

	"git.lolli.tech/lollipopkit/nano-db/api"
	"git.lolli.tech/lollipopkit/nano-db/consts"
	"git.lolli.tech/lollipopkit/nano-db/logger"
	"git.lolli.tech/lollipopkit/nano-db/utils"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	addr := flag.String("a", "0.0.0.0:3777", "specific the addr to listen")
	userName := flag.String("u", "", "generate the cookie with -n <username>")
	dbName := flag.String("d", "", "update acl rules with -d <dbname>")
	cacheLen := flag.Int("l", 100, "set the max length of cache")
	flag.Parse()

	initSalt()

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

	e.GET("/:db/:col", api.Files)
	e.DELETE("/:db/:col", api.DeleteCol)

	e.GET("/:db/:col/:file", api.Read)
	e.POST("/:db/:col/:file", api.Write)
	e.DELETE("/:db/:col/:file", api.Delete)

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
			println("this db already owned by other user")
			return
		}
		api.AclLock.RUnlock()
		println("you already owned this db")
		return
	}

	api.AclLock.RUnlock()
	api.AclLock.Lock()
	err := api.Acl.UpdateRule(*dbName, *userName)
	api.AclLock.Unlock()

	if err != nil {
		println("acl update rule: " + err.Error())
	} else {
		println("acl update rule: success")
	}
}

func initSalt() {
	if utils.IsExist(consts.SaltFile) {
		salt, err := ioutil.ReadFile(consts.SaltFile)
		if err != nil {
			println("[initSalt] ioutil.ReadFile(): " + err.Error())
			println("[initSalt] will use default salt")
			return
		}
		consts.CookieSalt = string(salt)
		return
	}
	runes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	salt := make([]rune, consts.SaltDefaultLen)
	for i := 0; i < consts.SaltDefaultLen; i++ {
		salt[i] = runes[rand.Intn(len(runes))]
	}
	ioutil.WriteFile(consts.SaltFile, []byte(string(salt)), consts.FilePermission)
	consts.CookieSalt = string(salt)
}
