package api

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"git.lolli.tech/lollipopkit/nano-db/consts"
	"git.lolli.tech/lollipopkit/nano-db/db"
	"git.lolli.tech/lollipopkit/nano-db/logger"
	"git.lolli.tech/lollipopkit/nano-db/model"
	"git.lolli.tech/lollipopkit/nano-db/utils"
	"github.com/labstack/echo"
)

var (
	cacher  = model.NewCacher(consts.CacherMaxLength * 100)
	Acl     = &model.ACL{}
	AclLock = &sync.RWMutex{}
)

const (
	pathFmt   = "%s/%s/%s"
	emptyPath = "[db] or [col] or [id] is empty"
)

func init() {
	go func() {
		AclLock.Lock()
		err := Acl.Load()
		AclLock.Unlock()

		if err != nil {
			panic(err)
		}
		time.Sleep(time.Minute)
	}()
}

func Exist(c echo.Context) error {
	dbName := c.Param("db")
	col := c.Param("col")
	id := c.Param("id")
	if dbName == "" || col == "" || id == "" {
		return resp(c, 520, emptyPath)
	}

	if !checkPermission(c, "api.Exist") {
		return resp(c, 403, "permission denied")
	}

	p := path(dbName, col, id)
	if err := verifyParams([]string{dbName, col, id}); err != nil {
		logger.W("[api.Exist] %s is not valid: %s\n", p, err.Error())
		return resp(c, 525, fmt.Sprintf("%s is not valid: %s", p, err.Error()))
	}

	exist := utils.IsExist(p)
	return resp(c, 200, exist)
}

func Read(c echo.Context) error {
	dbName := c.Param("db")
	col := c.Param("col")
	id := c.Param("id")
	if dbName == "" || col == "" || id == "" {
		return resp(c, 520, emptyPath)
	}

	if !checkPermission(c, "api.Read") {
		return resp(c, 403, "permission denied")
	}

	p := path(dbName, col, id)
	if err := verifyParams([]string{dbName, col, id}); err != nil {
		logger.W("[api.Write] %s is not valid: %s\n", p, err.Error())
		return resp(c, 525, fmt.Sprintf("%s is not valid: %s", p, err.Error()))
	}

	item, have := cacher.Get(p)
	if have {
		return resp(c, 200, item)
	}

	var content interface{}
	err := db.Read(p, &content)
	if err != nil {
		if err != db.ErrNoDocument {
			logger.E("[api.Read] db.Read(): %s\n", err.Error())
		}

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

	if !checkPermission(c, "api.Write") {
		return resp(c, 403, "permission denied")
	}

	p := path(dbName, col, id)
	if err := verifyParams([]string{dbName, col, id}); err != nil {
		logger.W("[api.Write] %s is not valid: %s\n", p, err.Error())
		return resp(c, 525, fmt.Sprintf("%s is not valid: %s", p, err.Error()))
	}

	var content interface{}
	err := c.Bind(&content)
	if err != nil {
		logger.E("[api.Write] c.Bind(): %s\n", err.Error())
		return resp(c, 522, "c.Bind(): "+err.Error())
	}

	err = os.MkdirAll(consts.DBDir+dbName+"/"+col, consts.FilePermission)
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

	if checkPermission(c, "api.Delete") {
		return resp(c, 403, "permission denied")
	}

	p := path(dbName, col, id)
	if err := verifyParams([]string{dbName, col, id}); err != nil {
		logger.W("[api.Write] %s is not valid: %s\n", p, err.Error())
		return resp(c, 525, fmt.Sprintf("%s is not valid: %s", p, err.Error()))
	}

	err := db.Delete(p)
	if err != nil {
		logger.E("[api.Delete] db.Delete(): %s\n", err.Error())
		return resp(c, 524, "db.Delete(): "+err.Error())
	}

	cacher.Delete(p)

	return resp(c, 200, "ok")
}

func IDs(c echo.Context) error {
	dbName := c.Param("db")
	col := c.Param("col")
	if dbName == "" || col == "" {
		return resp(c, 520, emptyPath)
	}

	if !checkPermission(c, "api.IDs") {
		return resp(c, 403, "permission denied")
	}

	p := path(dbName, col, "")
	ids, err := ioutil.ReadDir(p)
	if err != nil {
		logger.E("[api.IDs] ioutil.ReadDir(): %s\n", err.Error())
		return resp(c, 526, "ioutil.ReadDir(): "+err.Error())
	}

	var idsList []string
	for _, id := range ids {
		if !id.IsDir() {
			idsList = append(idsList, id.Name())
		}
	}

	return resp(c, 200, idsList)
}

func Cols(c echo.Context) error {
	dbName := c.Param("db")
	if dbName == "" {
		return resp(c, 520, emptyPath)
	}

	if !checkPermission(c, "api.Cols") {
		return resp(c, 403, "permission denied")
	}

	cols, err := ioutil.ReadDir(consts.DBDir + dbName)
	if err != nil {
		logger.E("[api.Cols] ioutil.ReadDir(): %s\n", err.Error())
		return resp(c, 527, "ioutil.ReadDir(): "+err.Error())
	}

	var colsList []string
	for _, col := range cols {
		if col.IsDir() {
			colsList = append(colsList, col.Name())
		}
	}

	return resp(c, 200, colsList)
}

func DeleteDB(c echo.Context) error {
	dbName := c.Param("db")
	if dbName == "" {
		return resp(c, 520, emptyPath)
	}

	if !checkPermission(c, "api.DeleteDB") {
		return resp(c, 403, "permission denied")
	}

	err := os.RemoveAll(consts.DBDir + dbName)
	if err != nil {
		logger.E("[api.DeleteDB] os.RemoveAll(): %s\n", err.Error())
		return resp(c, 528, "os.RemoveAll(): "+err.Error())
	}

	for _, path := range cacher.All() {
		p, ok := path.(string)
		if !ok {
			continue
		}
		if strings.HasPrefix(p, dbName+"/") {
			cacher.Delete(path)
		}
	}

	return resp(c, 200, "ok")
}

func DeleteCol(c echo.Context) error {
	dbName := c.Param("db")
	col := c.Param("col")
	if dbName == "" || col == "" {
		return resp(c, 520, emptyPath)
	}

	if !checkPermission(c, "api.DeleteCol") {
		return resp(c, 403, "permission denied")
	}

	err := os.RemoveAll(consts.DBDir + dbName + "/" + col)
	if err != nil {
		logger.E("[api.DeleteCol] os.RemoveAll(): %s\n", err.Error())
		return resp(c, 529, "os.RemoveAll(): "+err.Error())
	}

	for _, path := range cacher.All() {
		p, ok := path.(string)
		if !ok {
			continue
		}
		if strings.HasPrefix(p, dbName+"/"+col+"/") {
			cacher.Delete(path)
		}
	}

	return resp(c, 200, "ok")
}
