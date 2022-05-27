package api

import (
	"os"

	"github.com/LollipopKit/nano-db/consts"
	"github.com/LollipopKit/nano-db/db"
	"github.com/LollipopKit/nano-db/logger"
	"github.com/LollipopKit/nano-db/model"
	"github.com/labstack/echo"
)

var (
	cacher = model.NewCacher(consts.CacherMaxLength * 100)
)

const (
	pathFmt   = "%s/%s/%s"
	emptyPath = "[db] or [col] or [id] is empty"
)

func Read(c echo.Context) error {
	dbName := c.Param("db")
	col := c.Param("col")
	id := c.Param("id")
	if dbName == "" || col == "" || id == "" {
		return resp(c, 520, emptyPath)
	}

	loggedIn, userName := accountVerify(c)
	if !loggedIn {
		if userName != consts.AnonymousUser {
			logger.W("[api.Read] user %s is trying to read\n", userName)
		}
		return resp(c, 403, "permission denied")
	}

	if !verifyId(id) {
		logger.W("[api.Read] id %s is not valid\n", id)
		return resp(c, 525, "id is not valid")
	}

	p := path(dbName, col, id)

	item, have := cacher.Get(p)
	if have {
		return resp(c, 200, item)
	}

	var content interface{}
	err := db.Read(p, &content)
	if err != nil {
		logger.E("[api.Read] db.Read(): %s\n", err.Error())
		return resp(c, 521, "db.Read(): "+err.Error())
	}

	return resp(c, 200, content)
}

func Write(c echo.Context) error {
	dbName := c.Param("db")
	col := c.Param("col")
	id := c.Param("id")
	if dbName == "" || col == "" || id == "" {
		return resp(c, 520, emptyPath)
	}

	loggedIn, userName := accountVerify(c)
	if !loggedIn {
		if userName != consts.AnonymousUser {
			logger.W("[api.Write] user %s is trying to write\n", userName)
		}
		return resp(c, 403, "permission denied")
	}

	if !verifyId(id) {
		logger.W("[api.Write] id %s is not valid\n", id)
		return resp(c, 525, "id is not valid")
	}

	var content interface{}
	err := c.Bind(&content)
	if err != nil {
		logger.E("[api.Write] c.Bind(): %s\n", err.Error())
		return resp(c, 522, "c.Bind(): "+err.Error())
	}

	p := path(dbName, col, id)

	if err := os.MkdirAll(consts.DBDir + dbName + "/" + col, consts.FilePermission); err != nil {
		return err
	}

	err = db.Write(p, content)
	if err != nil {
		logger.E("[api.Write] db.Write(): %s\n", err.Error())
		return resp(c, 523, "db.Write(): "+err.Error())
	}

	cacher.Update(p, content)

	return resp(c, 200, nil)
}

func Delete(c echo.Context) error {
	dbName := c.Param("db")
	col := c.Param("col")
	id := c.Param("id")
	if dbName == "" || col == "" || id == "" {
		return resp(c, 520, emptyPath)
	}

	loggedIn, userName := accountVerify(c)
	if !loggedIn {
		if userName != consts.AnonymousUser {
			logger.W("[api.Delete] user %s is trying to delete\n", userName)
		}
		return resp(c, 403, "permission denied")
	}

	if !verifyId(id) {
		logger.W("[api.Delete] id %s is not valid\n", id)
		return resp(c, 525, "id is not valid")
	}

	p := path(dbName, col, id)

	err := db.Delete(p)
	if err != nil {
		logger.E("[api.Delete] db.Delete(): %s\n", err.Error())
		return resp(c, 524, "db.Delete(): "+err.Error())
	}

	cacher.Delete(p)

	return resp(c, 200, nil)
}
