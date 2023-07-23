package api

import (
	"io"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/lollipopkit/gommon/log"
	"github.com/lollipopkit/nano-db/cst"
)

const (
	statusFmt = "%d dbs, %d dirs in %s"
)

func Alive(c echo.Context) error {
	return c.NoContent(200)
}

func Read(c echo.Context) error {
	dbName := c.Param("db")
	dir := c.Param("dir")
	file := c.Param("file")

	p, err := path(dbName, dir, file)
	if err != nil {
		// 这里的错误可能是因为某些无权限的攻击引起的，所以不 log 记录
		// 同理，api.checkPermission 之前的错误都不 log 记录
		return c.String(cst.ECPath, err.Error())
	}

	if !checkPermission(c, "api.Read", dbName) {
		return permissionDenied(c)
	}

	return send(c, p)
}

func Write(c echo.Context) error {
	dbName := c.Param("db")
	dir := c.Param("dir")
	file := c.Param("file")

	p, err := path(dbName, dir, file)
	if err != nil {
		return c.String(cst.ECPath, err.Error())
	}

	if !checkPermission(c, "api.Write", dbName) {
		return permissionDenied(c)
	}

	err = os.MkdirAll(filepath.Dir(p), cst.FilePermission)
	if err != nil {
		log.Err("[api.Write] os.MkdirAll(): %s", err.Error())
		return c.String(cst.ECIO, "os.MkdirAll(): "+err.Error())
	}

	data, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Err("[api.Write] io.ReadAll(): %s", err.Error())
		return c.String(cst.ECIO, "io.ReadAll(): "+err.Error())
	}

	err = os.WriteFile(p, data, cst.FilePermission)
	if err != nil {
		log.Err("[api.Write] os.WriteFile(): %s", err.Error())
		return c.String(cst.ECIO, "os.WriteFile(): "+err.Error())
	}

	return c.NoContent(200)
}

func Delete(c echo.Context) error {
	dbName := c.Param("db")
	dir := c.Param("dir")
	file := c.Param("file")

	p, err := path(dbName, dir, file)
	if err != nil {
		return c.String(cst.ECIO, err.Error())
	}

	if !checkPermission(c, "api.Delete", dbName) {
		return permissionDenied(c)
	}

	err = os.Remove(p)
	if err != nil {
		log.Err("[api.Delete] os.Remove(): %s", err.Error())
		return c.String(cst.ECIO, "os.Remove(): "+err.Error())
	}

	return c.NoContent(200)
}

func ReadDir(c echo.Context) error {
	dbName := c.Param("db")
	dir := c.Param("dir")

	p, err := path(dbName, dir)
	if err != nil {
		return c.String(cst.ECPath, err.Error())
	}

	if !checkPermission(c, "api.ReadDir", dbName) {
		return permissionDenied(c)
	}

	files, err := os.ReadDir(p)
	if err != nil {
		log.Err("[api.ReadDir] os.ReadDir(): %s", err.Error())
		return c.String(cst.ECIO, "os.ReadDir(): "+err.Error())
	}

	var filesList []string
	for _, file := range files {
		if !file.IsDir() {
			filesList = append(filesList, file.Name())
		}
	}

	return c.JSON(200, filesList)
}

func ReadDB(c echo.Context) error {
	dbName := c.Param("db")

	p, err := path(dbName)
	if err != nil {
		return c.String(cst.ECPath, err.Error())
	}

	if !checkPermission(c, "api.ReadDB", dbName) {
		return permissionDenied(c)
	}

	dirs, err := os.ReadDir(p)
	if err != nil {
		log.Err("[api.ReadDB] os.ReadDir(): %s", err.Error())
		return c.String(cst.ECIO, "os.ReadDir(): "+err.Error())
	}

	var dirsList []string
	for _, dir := range dirs {
		if dir.IsDir() {
			dirsList = append(dirsList, dir.Name())
		}
	}

	return c.JSON(200, dirsList)
}

func DeleteDB(c echo.Context) error {
	dbName := c.Param("db")

	p, err := path(dbName)
	if err != nil {
		return c.String(cst.ECPath, err.Error())
	}

	if !checkPermission(c, "api.DeleteDB", dbName) {
		return permissionDenied(c)
	}

	err = os.RemoveAll(p)
	if err != nil {
		log.Err("[api.DeleteDB] os.RemoveAll(): %s", err.Error())
		return c.String(cst.ECIO, "os.RemoveAll(): "+err.Error())
	}

	return c.NoContent(200)
}

func DeleteDir(c echo.Context) error {
	dbName := c.Param("db")
	dir := c.Param("dir")

	p, err := path(dbName, dir)
	if err != nil {
		return c.String(cst.ECPath, err.Error())
	}

	if !checkPermission(c, "api.DeleteDir", dbName) {
		return permissionDenied(c)
	}

	err = os.RemoveAll(p)
	if err != nil {
		log.Err("[api.DeleteDir] os.RemoveAll(): %s", err.Error())
		return c.String(cst.ECIO, "os.RemoveAll(): "+err.Error())
	}

	return c.NoContent(200)
}
