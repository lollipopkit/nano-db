package api

import (
	"fmt"
	"io/ioutil"
	"time"

	"git.lolli.tech/LollipopKit/nano-db/consts"
	"git.lolli.tech/LollipopKit/nano-db/logger"
	"github.com/labstack/echo"
)

const (
	statusFmt = "%d dirs, %d files, %d cached items in %s"
)

func Home(c echo.Context) error {
	return resp(c, 200, "db alive")
}

func Status(c echo.Context) error {
	loggedIn, userName := accountVerify(c)
	if !loggedIn {
		if userName != consts.AnonymousUser {
			logger.W("[api.Status] user %s is trying to get\n", userName)
		}
		return resp(c, 403, "permission denied")
	}

	time1 := time.Now()

	dirs, err := ioutil.ReadDir(consts.DBDir)
	if err != nil {
		return resp(c, 525, "ioutil.ReadDir(): "+err.Error())
	}

	dirNames := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		if dir.IsDir() {
			dirNames = append(dirNames, dir.Name())
		}
	}

	filesCount := 0
	for _, dirName := range dirNames {
		files, err := ioutil.ReadDir(consts.DBDir + dirName)
		if err != nil {
			return resp(c, 525, "ioutil.ReadDir(): "+err.Error())
		}

		filesCount += len(files)
	}

	cacherLen := cacher.Len()

	time2 := time.Now()

	return resp(c, 200, fmt.Sprintf(statusFmt, len(dirNames), filesCount, cacherLen, time2.Sub(time1).String()))
}
