package api

import (
	"os"
	"sync"

	"git.lolli.tech/LollipopKit/nano-db/consts"
	"git.lolli.tech/LollipopKit/nano-db/db"
	"git.lolli.tech/LollipopKit/nano-db/logger"
	"git.lolli.tech/LollipopKit/nano-db/model"
	"github.com/labstack/echo"
)

var (
	cacher = model.NewCacher(consts.CacherMaxLength * 100)
	acl = &model.ACL{}
	aclLock = &sync.RWMutex{}
)

const (
	pathFmt   = "%s/%s/%s"
	emptyPath = "[db] or [col] or [id] is empty"
)

func init() {
	aclLock.Lock()
	err := acl.Load()
	aclLock.Unlock()

	if err != nil {
		panic(err)
	}
}

func Init(c echo.Context) error {
	dbName := c.Param("db")
	if dbName == "" {
		return resp(c, 520, "dbName is empty")
	}

	loggedIn, userName := accountVerify(c)
	if !loggedIn {
		if userName != consts.AnonymousUser {
			logger.W("[api.Init] user %s is trying to init\n", userName)
		}
		return resp(c, 403, "permission denied")
	}

	aclLock.RLock()
	if acl.HaveDB(dbName) {
		if !acl.Can(dbName, userName) {
			aclLock.RUnlock()
			return resp(c, 403, "this db already initialized by other user")
		}
		aclLock.RUnlock()
		return resp(c, 200, "you already initialized this db")
	}

	aclLock.RUnlock()
	aclLock.Lock()
	err := acl.UpdateRule(dbName, userName)
	aclLock.Unlock()

	if err != nil {
		logger.E("[api.Init] acl.UpdateRule(): %s\n", err.Error())
		return resp(c, 526, "acl.UpdateRule(): "+err.Error())
	}

	return resp(c, 200, "ok")
}

func Read(c echo.Context) error {
	dbName := c.Param("db")
	col := c.Param("col")
	id := c.Param("id")
	if dbName == "" || col == "" || id == "" {
		return resp(c, 520, emptyPath)
	}

	loggedIn, userName := accountVerify(c)
	aclLock.RLock()
	if !loggedIn || !acl.Can(dbName, userName) {
		if userName != consts.AnonymousUser {
			logger.W("[api.Read] user %s is trying to read %s\n", userName, dbName)
		}
		aclLock.RUnlock()
		return resp(c, 403, "permission denied")
	}
	aclLock.RUnlock()

	if !verifyParams([]string{dbName, col, id}) {
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

	cacher.Update(p, content)

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
	aclLock.RLock()
	if !loggedIn || !acl.Can(dbName, userName) {
		if userName != consts.AnonymousUser {
			logger.W("[api.Write] user %s is trying to write %s\n", userName, dbName)
		}
		aclLock.RUnlock()
		return resp(c, 403, "permission denied")
	}
	aclLock.RUnlock()

	if !verifyParams([]string{dbName, col, id}) {
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

	err = os.MkdirAll(consts.DBDir + dbName + "/" + col, consts.FilePermission)
	if err != nil {
		logger.E("[api.Write] os.MkdirAll(): %s\n", err.Error())
		return resp(c, 523, "os.MkdirAll(): "+err.Error())
	}

	err = db.Write(p, content)
	if err != nil {
		logger.E("[api.Write] db.Write(): %s\n", err.Error())
		return resp(c, 523, "db.Write(): "+err.Error())
	}

	cacher.Update(p, content)

	return resp(c, 200, "ok")
}

func Delete(c echo.Context) error {
	dbName := c.Param("db")
	col := c.Param("col")
	id := c.Param("id")
	if dbName == "" || col == "" || id == "" {
		return resp(c, 520, emptyPath)
	}

	loggedIn, userName := accountVerify(c)
	aclLock.RLock()
	if !loggedIn || !acl.Can(dbName, userName) {
		if userName != consts.AnonymousUser {
			logger.W("[api.Delete] user %s is trying to delete %s/%s/%s\n", userName, dbName, col, id)
		}
		aclLock.RUnlock()
		return resp(c, 403, "permission denied")
	}
	aclLock.RUnlock()

	if !verifyParams([]string{dbName, col, id}) {
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

	return resp(c, 200, "ok")
}
