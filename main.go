package main

import (
	"flag"
	"math/rand"
	"os"
	"regexp"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lollipopkit/gommon/log"
	"github.com/lollipopkit/nano-db/api"
	. "github.com/lollipopkit/nano-db/cfg"
)

func main() {
	parseCli()
	if err := startWeb(); err != nil {
		log.Err(err.Error())
	}
}

func parseCli() {
	token := flag.String("t", "", "set token with -u <token>")
	dbName := flag.String("d", "", "update acl rules with -d <dbname>")
	flag.Parse()

	// generate cookie & update acl rules
	if *dbName != "" {
		if *token == "" {
			*token = genToken()
			log.Info("generated token: %s", *token)
		}
		UpdateAcl(*token, *dbName)
		os.Exit(0)
	}
}

func startWeb() error {
	e := echo.New()

	if App.Log.Enable {
		if App.Log.Format == "" {
			e.Use(middleware.Logger())
		} else {
			skipRegList := make([]*regexp.Regexp, 0, len(App.Log.SkipRegExp))
			for _, reg := range App.Log.SkipRegExp {
				skipRegList = append(skipRegList, regexp.MustCompile(reg))
			}

			e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
				Format: App.Log.Format,
				Skipper: func(context echo.Context) bool {
					url := context.Request().URL.Path
					for _, reg := range skipRegList {
						if reg.MatchString(url) {
							return true
						}
					}
					return false
				},
			}))
		}
	}
	e.Use(api.RateLimiter)
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: App.Security.CORSList,
	}))
	e.Use(middleware.BodyLimit(App.Security.BodyLimit))

	// Routes
	e.HEAD("/", api.Alive)

	e.GET("/:db", api.ReadDB, api.CheckPathAndPerm(1))
	e.DELETE("/:db", api.DeleteDB, api.CheckPathAndPerm(1))

	e.GET("/:db/:dir", api.ReadDir, api.CheckPathAndPerm(2))
	e.DELETE("/:db/:dir", api.DeleteDir, api.CheckPathAndPerm(2))

	e.GET("/:db/:dir/:file", api.Read, api.CheckPathAndPerm(3))
	e.POST("/:db/:dir/:file", api.Write, api.CheckPathAndPerm(3))
	e.DELETE("/:db/:dir/:file", api.Delete, api.CheckPathAndPerm(3))

	e.HTTPErrorHandler = api.HandleErr

	// Start server
	e.HideBanner = true
	return e.Start(App.Addr)
}

func genToken() string {
	runes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	salt := make([]rune, App.Security.TokenLen)
	for i := 0; i < App.Security.TokenLen; i++ {
		salt[i] = runes[rand.Intn(len(runes))]
	}
	return string(salt)
}
