package api

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	glc "git.lolli.tech/lollipopkit/go_lru_cacher"
	"git.lolli.tech/lollipopkit/nano-db/consts"
	"git.lolli.tech/lollipopkit/nano-db/db"
	"git.lolli.tech/lollipopkit/nano-db/logger"
	"git.lolli.tech/lollipopkit/nano-db/model"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/tidwall/gjson"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary

	cacher = glc.NewCacher(consts.CacherMaxLength * 100)
	// record ip's failed times: map[string]int -> {"ip": times}
	banIP = glc.NewCacher(consts.CacherMaxLength)

	Acl     = &model.ACL{}
	AclLock = &sync.RWMutex{}
)

const (
	pathFmt        = "%s/%s/%s"
	emptyPath      = "[db] or [dir] or [file] is empty"
	emptyGJsonPath = "gjson path is empty"
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

func checkIP(c echo.Context) (int, error) {
	v, ok := banIP.Get(c.RealIP())
	if ok {
		times, ok := v.(int)
		if ok {
			if times >= consts.MaxIPFailedTimes {
				return times, resp(c, 531, "ip is blocked")
			}
		}
		return times, resp(c, 532, "banned ip list convert error")
	}
	return 0, nil
}

func Read(c echo.Context) error {
	banTimes, err := checkIP(c)
	if err != nil {
		return err
	}

	dbName := c.Param("db")
	dir := c.Param("dir")
	file := c.Param("file")
	if dbName == "" || dir == "" || file == "" {
		return resp(c, 520, emptyPath)
	}

	if !checkPermission(c, "api.Read") {
		banIP.Set(c.RealIP(), banTimes+1)
		return permissionDenied(c)
	}

	p := path(dbName, dir, file)
	if err := verifyParams([]string{dbName, dir, file}); err != nil {
		logger.W("[api.Write] %s is not valid: %s\n", p, err.Error())
		return resp(c, 525, fmt.Sprintf("%s is not valid: %s", p, err.Error()))
	}

	item, have := cacher.Get(p)
	if have {
		return resp(c, 200, item)
	}

	var content any
	err = db.Read(p, &content)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			logger.E("[api.Read] db.Read(): %s\n", err.Error())
		}

		return resp(c, 521, "db.Read(): "+err.Error())
	}

	cacher.Set(p, content)

	return resp(c, 200, content)
}

func Write(c echo.Context) error {
	banTimes, err := checkIP(c)
	if err != nil {
		return err
	}

	dbName := c.Param("db")
	dir := c.Param("dir")
	file := c.Param("file")
	if dbName == "" || dir == "" || file == "" {
		return resp(c, 520, emptyPath)
	}

	if !checkPermission(c, "api.Write") {
		banIP.Set(c.RealIP(), banTimes+1)
		return permissionDenied(c)
	}

	p := path(dbName, dir, file)
	if err := verifyParams([]string{dbName, dir, file}); err != nil {
		logger.W("[api.Write] %s is not valid: %s\n", p, err.Error())
		return resp(c, 525, fmt.Sprintf("%s is not valid: %s", p, err.Error()))
	}

	var content any
	err = c.Bind(&content)
	if err != nil {
		logger.E("[api.Write] c.Bind(): %s\n", err.Error())
		return resp(c, 522, "c.Bind(): "+err.Error())
	}

	err = os.MkdirAll(consts.DBDir+dbName+"/"+dir, consts.FilePermission)
	if err != nil {
		logger.E("[api.Write] os.MkdirAll(): %s\n", err.Error())
		return resp(c, 523, "os.MkdirAll(): "+err.Error())
	}

	err = db.Write(p, content)
	if err != nil {
		logger.E("[api.Write] db.Write(): %s\n", err.Error())
		return resp(c, 523, "db.Write(): "+err.Error())
	}

	cacher.Set(p, content)

	return ok(c)
}

func Delete(c echo.Context) error {
	banTimes, err := checkIP(c)
	if err != nil {
		return err
	}

	dbName := c.Param("db")
	dir := c.Param("dir")
	file := c.Param("file")
	if dbName == "" || dir == "" || file == "" {
		return resp(c, 520, emptyPath)
	}

	if !checkPermission(c, "api.Delete") {
		banIP.Set(c.RealIP(), banTimes+1)
		return permissionDenied(c)
	}

	p := path(dbName, dir, file)
	if err := verifyParams([]string{dbName, dir, file}); err != nil {
		logger.W("[api.Write] %s is not valid: %s\n", p, err.Error())
		return resp(c, 525, fmt.Sprintf("%s is not valid: %s", p, err.Error()))
	}

	err = db.Delete(p)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			logger.E("[api.Delete] db.Delete(): %s\n", err.Error())
		}
		return resp(c, 524, "db.Delete(): "+err.Error())
	}

	cacher.Delete(p)

	return ok(c)
}

func Files(c echo.Context) error {
	banTimes, err := checkIP(c)
	if err != nil {
		return err
	}

	dbName := c.Param("db")
	dir := c.Param("dir")
	if dbName == "" || dir == "" {
		return resp(c, 520, emptyPath)
	}

	if !checkPermission(c, "api.Files") {
		banIP.Set(c.RealIP(), banTimes+1)
		return permissionDenied(c)
	}

	p := consts.DBDir + path(dbName, dir, "")
	files, err := ioutil.ReadDir(p)
	if err != nil {
		logger.E("[api.IDs] ioutil.ReadDir(): %s\n", err.Error())
		return resp(c, 526, "ioutil.ReadDir(): "+err.Error())
	}

	var filesList []string
	for _, file := range files {
		if !file.IsDir() {
			filesList = append(filesList, file.Name())
		}
	}

	return resp(c, 200, filesList)
}

func Dirs(c echo.Context) error {
	banTimes, err := checkIP(c)
	if err != nil {
		return err
	}

	dbName := c.Param("db")
	if dbName == "" {
		return resp(c, 520, emptyPath)
	}

	if !checkPermission(c, "api.Dirs") {
		banIP.Set(c.RealIP(), banTimes+1)
		return permissionDenied(c)
	}

	dirs, err := ioutil.ReadDir(consts.DBDir + dbName)
	if err != nil {
		logger.E("[api.Dirs] ioutil.ReadDir(): %s\n", err.Error())
		return resp(c, 527, "ioutil.ReadDir(): "+err.Error())
	}

	var dirsList []string
	for _, dir := range dirs {
		if dir.IsDir() {
			dirsList = append(dirsList, dir.Name())
		}
	}

	return resp(c, 200, dirsList)
}

func DeleteDB(c echo.Context) error {
	banTimes, err := checkIP(c)
	if err != nil {
		return err
	}

	dbName := c.Param("db")
	if dbName == "" {
		return resp(c, 520, emptyPath)
	}

	if !checkPermission(c, "api.DeleteDB") {
		banIP.Set(c.RealIP(), banTimes+1)
		return permissionDenied(c)
	}

	err = os.RemoveAll(consts.DBDir + dbName)
	if err != nil {
		logger.E("[api.DeleteDB] os.RemoveAll(): %s\n", err.Error())
		return resp(c, 528, "os.RemoveAll(): "+err.Error())
	}

	for _, path := range cacher.Values() {
		p, ok := path.(string)
		if !ok {
			continue
		}
		if strings.HasPrefix(p, dbName+"/") {
			cacher.Delete(path)
		}
	}

	return ok(c)
}

func DeleteCol(c echo.Context) error {
	banTimes, err := checkIP(c)
	if err != nil {
		return err
	}

	dbName := c.Param("db")
	dir := c.Param("dir")
	if dbName == "" || dir == "" {
		return resp(c, 520, emptyPath)
	}

	if !checkPermission(c, "api.DeleteCol") {
		banIP.Set(c.RealIP(), banTimes+1)
		return permissionDenied(c)
	}

	err = os.RemoveAll(consts.DBDir + dbName + "/" + dir)
	if err != nil {
		logger.E("[api.DeleteCol] os.RemoveAll(): %s\n", err.Error())
		return resp(c, 529, "os.RemoveAll(): "+err.Error())
	}

	for _, path := range cacher.Values() {
		p, ok := path.(string)
		if !ok {
			continue
		}
		if strings.HasPrefix(p, dbName+"/"+dir+"/") {
			cacher.Delete(path)
		}
	}

	return ok(c)
}

func SearchInDir(c echo.Context) error {
	banTimes, err := checkIP(c)
	if err != nil {
		return err
	}

	dbName := c.Param("db")
	dir := c.Param("dir")
	if dbName == "" || dir == "" {
		return resp(c, 520, emptyPath)
	}

	var searchReq model.SearchReq
	err = c.Bind(&searchReq)
	if err != nil {
		logger.E("[api.Search] c.Bind(): %s\n", err.Error())
		return resp(c, 530, "c.Bind(): "+err.Error())
	}

	gjsonPath := searchReq.Path
	if gjsonPath == "" {
		return resp(c, 521, emptyGJsonPath)
	}
	valueRegex := searchReq.Regex

	if !checkPermission(c, "api.Search") {
		banIP.Set(c.RealIP(), banTimes+1)
		return permissionDenied(c)
	}

	p := consts.DBDir + path(dbName, dir, "")
	files, err := ioutil.ReadDir(p)
	if err != nil {
		logger.E("[api.Search] ioutil.ReadDir(): %s\n", err.Error())
		return resp(c, 530, "ioutil.ReadDir(): "+err.Error())
	}

	var results []any
	for _, file := range files {
		var data []byte
		var err error
		var ok bool
		var d any

		d, ok = cacher.Get(path(dbName, dir, file.Name()))
		if ok {
			data, err = json.Marshal(d)
			if err != nil {
				logger.E("[api.Search] json.Marshal(): %s\n", err.Error())
				continue
			}
		} else {
			data, err = ioutil.ReadFile(p + file.Name())
			if err != nil {
				logger.E("[api.Search] ioutil.ReadFile(): %s\n", err.Error())
				continue
			}
			err = json.Unmarshal(data, &d)
			if err != nil {
				logger.E("[api.Search] json.Unmarshal(): %s\n", err.Error())
				continue
			}
		}

		result := gjson.GetBytes(data, gjsonPath)
		if result.Exists() {
			if ok, err := regexp.MatchString(valueRegex, result.Raw); (err == nil && ok) || valueRegex == "" {
				results = append(results, d)
			}
		}
	}

	return resp(c, 200, results)
}

func SearchInDB(c echo.Context) error {
	banTimes, err := checkIP(c)
	if err != nil {
		return err
	}

	dbName := c.Param("db")
	if dbName == "" {
		return resp(c, 520, emptyPath)
	}

	var searchReq model.SearchReq
	err = c.Bind(&searchReq)
	if err != nil {
		logger.E("[api.Search] c.Bind(): %s\n", err.Error())
		return resp(c, 530, "c.Bind(): "+err.Error())
	}

	gjsonPath := searchReq.Path
	if gjsonPath == "" {
		return resp(c, 521, emptyGJsonPath)
	}
	valueRegex := searchReq.Regex

	if !checkPermission(c, "api.Search") {
		banIP.Set(c.RealIP(), banTimes+1)
		return permissionDenied(c)
	}

	p := consts.DBDir + dbName + "/"
	dirs, err := ioutil.ReadDir(p)
	if err != nil {
		logger.E("[api.Search] ioutil.ReadDir(): %s\n", err.Error())
		return resp(c, 530, "ioutil.ReadDir(): "+err.Error())
	}

	var results []any
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		files, err := ioutil.ReadDir(p + dir.Name() + "/")
		if err != nil {
			logger.E("[api.Search] ioutil.ReadDir(): %s\n", err.Error())
			continue
		}
		for _, file := range files {
			var data []byte
			var err error
			var ok bool
			var d any

			d, ok = cacher.Get(path(dbName, dir.Name(), file.Name()))
			if ok {
				data, err = json.Marshal(d)
				if err != nil {
					logger.E("[api.Search] json.Marshal(): %s\n", err.Error())
					continue
				}
			} else {
				data, err = ioutil.ReadFile(p + dir.Name() + "/" + file.Name())
				if err != nil {
					logger.E("[api.Search] ioutil.ReadFile(): %s\n", err.Error())
					continue
				}
				err = json.Unmarshal(data, &d)
				if err != nil {
					logger.E("[api.Search] json.Unmarshal(): %s\n", err.Error())
					continue
				}
			}

			result := gjson.GetBytes(data, gjsonPath)
			if result.Exists() {
				if ok, err := regexp.MatchString(valueRegex, result.Raw); (err == nil && ok) || valueRegex == "" {
					results = append(results, d)
				}
			}
		}
	}
	
	return resp(c, 200, results)
}
