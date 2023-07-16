package api

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	glc "github.com/lollipopkit/go-lru-cacher"
	"github.com/lollipopkit/gommon/log"
	. "github.com/lollipopkit/nano-db/cfg"
	"github.com/lollipopkit/nano-db/consts"
	"github.com/lollipopkit/nano-db/db"
	. "github.com/lollipopkit/nano-db/json"
	"github.com/tidwall/gjson"
)

var (
	_duration    = time.Hour
	dbDataCacher *glc.PartedCacher
)

const (
	emptyPath      = "[db] or [dir] or [file] is empty"
	emptyGJsonPath = "gjson path is empty"
)

func init() {
	dbDataCacher = glc.NewPartedElapsedCacher(
		Cfg.Cache.MaxSize*100,
		Cfg.Cache.ActiveRate,
		_duration,
		_duration*24,
	)
}

func Read(c echo.Context) error {
	dbName := c.Param("db")
	dir := c.Param("dir")
	file := c.Param("file")
	if dbName == "" || dir == "" || file == "" {
		return resp(c, 520, emptyPath)
	}

	p := path(dbName, dir, file)
	if !checkPermission(c, "api.Read", dbName, p) {
		return permissionDenied(c)
	}

	if err := verifyParams([]string{dbName, dir, file}); err != nil {
		log.Warn("[api.Write] %s is not valid: %s", p, err.Error())
		return resp(c, 525, fmt.Sprintf("%s is not valid: %s", p, err.Error()))
	}

	item, have := dbDataCacher.Get(p)
	if have {
		return resp(c, 200, item)
	}

	var content any
	err := db.Read(p, &content)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Err("[api.Read] db.Read(): %s", err.Error())
		}

		return resp(c, 521, "db.Read(): "+err.Error())
	}

	dbDataCacher.Set(p, content)

	return resp(c, 200, content)
}

func Write(c echo.Context) error {
	dbName := c.Param("db")
	dir := c.Param("dir")
	file := c.Param("file")
	if dbName == "" || dir == "" || file == "" {
		return resp(c, 520, emptyPath)
	}

	p := path(dbName, dir, file)
	if !checkPermission(c, "api.Write", dbName, p) {
		return permissionDenied(c)
	}

	if err := verifyParams([]string{dbName, dir, file}); err != nil {
		log.Warn("[api.Write] %s is not valid: %s", p, err.Error())
		return resp(c, 525, fmt.Sprintf("%s is not valid: %s", p, err.Error()))
	}

	var content any
	err := c.Bind(&content)
	if err != nil {
		log.Err("[api.Write] c.Bind(): %s", err.Error())
		return resp(c, 522, "c.Bind(): "+err.Error())
	}

	err = os.MkdirAll(consts.DBDir+dbName+"/"+dir, consts.FilePermission)
	if err != nil {
		log.Err("[api.Write] os.MkdirAll(): %s", err.Error())
		return resp(c, 523, "os.MkdirAll(): "+err.Error())
	}

	err = db.Write(p, content)
	if err != nil {
		log.Err("[api.Write] db.Write(): %s", err.Error())
		return resp(c, 523, "db.Write(): "+err.Error())
	}

	dbDataCacher.Set(p, content)

	return ok(c)
}

func Delete(c echo.Context) error {
	dbName := c.Param("db")
	dir := c.Param("dir")
	file := c.Param("file")
	if dbName == "" || dir == "" || file == "" {
		return resp(c, 520, emptyPath)
	}

	p := path(dbName, dir, file)
	if !checkPermission(c, "api.Delete", dbName, p) {
		return permissionDenied(c)
	}

	if err := verifyParams([]string{dbName, dir, file}); err != nil {
		log.Warn("[api.Write] %s is not valid: %s", p, err.Error())
		return resp(c, 525, fmt.Sprintf("%s is not valid: %s", p, err.Error()))
	}

	err := db.Delete(p)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Err("[api.Delete] db.Delete(): %s", err.Error())
		}
		return resp(c, 524, "db.Delete(): "+err.Error())
	}

	dbDataCacher.Delete(p)

	return ok(c)
}

func Files(c echo.Context) error {
	dbName := c.Param("db")
	dir := c.Param("dir")
	if dbName == "" || dir == "" {
		return resp(c, 520, emptyPath)
	}

	p := consts.DBDir + path(dbName, dir, "")
	if !checkPermission(c, "api.Files", dbName, p) {
		return permissionDenied(c)
	}

	files, err := os.ReadDir(p)
	if err != nil {
		log.Err("[api.IDs] os.ReadDir(): %s", err.Error())
		return resp(c, 526, "os.ReadDir(): "+err.Error())
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
	dbName := c.Param("db")
	if dbName == "" {
		return resp(c, 520, emptyPath)
	}

	if !checkPermission(c, "api.Dirs", dbName, dbName) {
		return permissionDenied(c)
	}

	dirs, err := os.ReadDir(consts.DBDir + dbName)
	if err != nil {
		log.Err("[api.Dirs] os.ReadDir(): %s", err.Error())
		return resp(c, 527, "os.ReadDir(): "+err.Error())
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
	dbName := c.Param("db")
	if dbName == "" {
		return resp(c, 520, emptyPath)
	}

	if !checkPermission(c, "api.DeleteDB", dbName, dbName) {
		return permissionDenied(c)
	}

	err := os.RemoveAll(consts.DBDir + dbName)
	if err != nil {
		log.Err("[api.DeleteDB] os.RemoveAll(): %s", err.Error())
		return resp(c, 528, "os.RemoveAll(): "+err.Error())
	}

	for _, path := range dbDataCacher.Values() {
		p, ok := path.(string)
		if !ok {
			continue
		}
		if strings.HasPrefix(p, dbName+"/") {
			dbDataCacher.Delete(path)
		}
	}

	return ok(c)
}

func DeleteDir(c echo.Context) error {
	dbName := c.Param("db")
	dir := c.Param("dir")
	if dbName == "" || dir == "" {
		return resp(c, 520, emptyPath)
	}

	if !checkPermission(c, "api.DeleteDir", dbName, dbName) {
		return permissionDenied(c)
	}

	err := os.RemoveAll(consts.DBDir + dbName + "/" + dir)
	if err != nil {
		log.Err("[api.DeleteDir] os.RemoveAll(): %s", err.Error())
		return resp(c, 529, "os.RemoveAll(): "+err.Error())
	}

	for _, path := range dbDataCacher.Values() {
		p, ok := path.(string)
		if !ok {
			continue
		}
		if strings.HasPrefix(p, dbName+"/"+dir+"/") {
			dbDataCacher.Delete(path)
		}
	}

	return ok(c)
}

func SearchDir(c echo.Context) error {
	dbName := c.Param("db")
	dir := c.Param("dir")
	if dbName == "" || dir == "" {
		return resp(c, 520, emptyPath)
	}

	p := consts.DBDir + path(dbName, dir, "")
	if !checkPermission(c, "api.SearchDir", dbName, p) {
		return permissionDenied(c)
	}

	searchReq := new(SearchReq)
	err := c.Bind(searchReq)
	if err != nil {
		log.Err("[api.SearchDir] c.Bind(): %s", err.Error())
		return resp(c, 530, "c.Bind(): "+err.Error())
	}

	if searchReq.Path == "" {
		return resp(c, 521, emptyGJsonPath)
	}

	files, err := os.ReadDir(p)
	if err != nil {
		log.Err("[api.SearchDir] os.ReadDir(): %s", err.Error())
		return resp(c, 530, "os.ReadDir(): "+err.Error())
	}

	var results []any
	for _, file := range files {
		var data []byte
		var err error
		var ok bool
		var d any

		d, ok = dbDataCacher.Get(path(dbName, dir, file.Name()))
		if ok {
			data, err = Json.Marshal(d)
			if err != nil {
				log.Err("[api.SearchDir] JsonMarshal(): %s", err.Error())
				continue
			}
		} else {
			data, err = os.ReadFile(p + file.Name())
			if err != nil {
				log.Err("[api.SearchDir] os.ReadFile(): %s", err.Error())
				continue
			}
			err = Json.Unmarshal(data, &d)
			if err != nil {
				log.Err("[api.SearchDir] JsonUnmarshal(): %s", err.Error())
				continue
			}
		}

		result := gjson.GetBytes(data, searchReq.Path)
		if result.Exists() {
			if searchReq.Regex == "" {
				results = append(results, d)
				continue
			}
			ok, err := regexp.MatchString(searchReq.Regex, result.Raw)
			if err == nil && ok {
				results = append(results, d)
			}
		}
	}

	return resp(c, 200, results)
}

func SearchDB(c echo.Context) error {
	dbName := c.Param("db")
	if dbName == "" {
		return resp(c, 520, emptyPath)
	}

	if !checkPermission(c, "api.SearchDB", dbName, dbName) {
		return permissionDenied(c)
	}

	searchReq := new(SearchReq)
	err := c.Bind(searchReq)
	if err != nil {
		log.Err("[api.SearchDB] c.Bind(): %s", err.Error())
		return resp(c, 530, "c.Bind(): "+err.Error())
	}

	if searchReq.Path == "" {
		return resp(c, 521, emptyGJsonPath)
	}

	p := consts.DBDir + dbName + "/"
	dirs, err := os.ReadDir(p)
	if err != nil {
		log.Err("[api.SearchDB] os.ReadDir(): %s", err.Error())
		return resp(c, 530, "os.ReadDir(): "+err.Error())
	}

	var results []any
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		files, err := os.ReadDir(p + dir.Name() + "/")
		if err != nil {
			log.Err("[api.SearchDB] os.ReadDir(): %s", err.Error())
			continue
		}
		for _, file := range files {
			var data []byte
			var err error
			var ok bool
			var d any

			d, ok = dbDataCacher.Get(path(dbName, dir.Name(), file.Name()))
			if ok {
				data, err = Json.Marshal(d)
				if err != nil {
					log.Err("[api.SearchDB] JsonMarshal(): %s", err.Error())
					continue
				}
			} else {
				data, err = os.ReadFile(p + dir.Name() + "/" + file.Name())
				if err != nil {
					log.Err("[api.SearchDB] os.ReadFile(): %s", err.Error())
					continue
				}
				err = Json.Unmarshal(data, &d)
				if err != nil {
					log.Err("[api.SearchDB] JsonUnmarshal(): %s", err.Error())
					continue
				}
			}

			result := gjson.GetBytes(data, searchReq.Path)
			if result.Exists() {
				if searchReq.Regex == "" {
					results = append(results, d)
					continue
				}
				ok, err := regexp.MatchString(searchReq.Regex, result.Raw)
				if err == nil && ok {
					results = append(results, d)
				}
			}
		}
	}

	return resp(c, 200, results)
}
