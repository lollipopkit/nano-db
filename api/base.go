package api

import (
	"fmt"
	"os"
	"time"

	. "git.lolli.tech/lollipopkit/nano-db/acl"
	"git.lolli.tech/lollipopkit/nano-db/consts"
	"git.lolli.tech/lollipopkit/nano-db/logger"
	"github.com/labstack/echo/v4"
)

const (
	statusFmt = "%d dbs, %d dirs, %d cached items in %s"
)

func Alive(c echo.Context) error {
	return c.NoContent(200)
}

func Status(c echo.Context) error {
	loggedIn, userName := accountVerify(c)
	if !loggedIn {
		if userName != consts.AnonymousUser {
			logger.W("[api.Status] user %s is trying to get\n", userName)
		}
		return permissionDenied(c)
	}

	time1 := time.Now()

	dirs, err := os.ReadDir(consts.DBDir)
	if err != nil {
		return resp(c, 525, "os.ReadDir(): "+err.Error())
	}

	dirNames := make([]string, 0, len(dirs))
	for _, d := range dirs {
		dbName := d.Name()
		if !d.IsDir() {
			logger.W("[api.Status] %s is not a dir\n", dbName)
			continue
		}
		if Acl.Can(dbName, userName) {
			dirNames = append(dirNames, dbName)
		}
	}

	filesCount := 0
	for _, dirName := range dirNames {
		files, err := os.ReadDir(consts.DBDir + dirName)
		if err != nil {
			return resp(c, 525, "os.ReadDir(): "+err.Error())
		}

		filesCount += len(files)
	}

	cacherLen := cacher.Len()

	time2 := time.Now()

	return resp(c, 200, fmt.Sprintf(statusFmt, len(dirNames), filesCount, cacherLen, time2.Sub(time1).String()))
}
