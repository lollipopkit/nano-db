package main

import (
	"flag"
	"os"
	"regexp"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lollipopkit/gommon/term"
	"github.com/lollipopkit/nano-db/api"
	. "github.com/lollipopkit/nano-db/cfg"
)

func main() {
	parseCli()
	if err := startWeb(); err != nil {
		term.Err(err.Error())
	}
}

func parseCli() {
	userName := flag.String("u", "", "generate the cookie with -u <username>")
	dbName := flag.String("d", "", "update acl rules with -d <dbname>")
	flag.Parse()

	// generate cookie & update acl rules
	if *userName != "" {
		if *dbName == "" {
			term.Info(api.GenCookie(*userName))
		} else {
			UpdateAcl(userName, dbName)
		}
		os.Exit(0)
	}
}

func startWeb() error {
	e := echo.New()

	if Cfg.Log.Enable {
		if Cfg.Log.Format == "" {
			e.Use(middleware.Logger())
		} else {
			skipRegList := make([]regexp.Regexp, 0, len(Cfg.Log.SkipRegExp))
			for _, reg := range Cfg.Log.SkipRegExp {
				skipRegList = append(skipRegList, *regexp.MustCompile(reg))
			}

			e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
				Format:  Cfg.Log.Format,
				Skipper: func (context echo.Context) bool {
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
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(Cfg.Security.RateLimit)))
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: Cfg.Security.CORSList,
	}))
	e.Use(middleware.BodyLimit(Cfg.Security.BodyLimit))

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
	return e.Start(Cfg.Addr)
}
