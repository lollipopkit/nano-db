package api

import (
	"fmt"
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

	p, err := checkAndJoinPath(dbName, dir, file)
	if err != nil {
		// 这里的错误可能是因为某些无权限的攻击引起的，所以不 log 记录
		// 同理，api.checkPermission 之前的错误都不 log 记录
		return c.String(cePath, err.Error())
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

	p, err := checkAndJoinPath(dbName, dir, file)
	if err != nil {
		return c.String(cePath, err.Error())
	}

	if !checkPermission(c, "api.Write", dbName) {
		return permissionDenied(c)
	}

	err = os.MkdirAll(filepath.Dir(p), cst.FilePermission)
	if err != nil {
		errStr := fmt.Sprintf("[api.Write] os.MkdirAll(): %s", err.Error())
		log.Err(errStr)
		return c.String(ceIO, errStr)
	}

	data, err := io.ReadAll(c.Request().Body)
	if err != nil {
		errStr := fmt.Sprintf("[api.Write] io.ReadAll(): %s", err.Error())
		log.Err(errStr)
		return c.String(ceIO, errStr)
	}

	err = os.WriteFile(p, data, cst.FilePermission)
	if err != nil {
		errStr := fmt.Sprintf("[api.Write] os.WriteFile(): %s", err.Error())
		log.Err(errStr)
		return c.String(ceIO, errStr)
	}

	return c.NoContent(200)
}

func Delete(c echo.Context) error {
	dbName := c.Param("db")
	dir := c.Param("dir")
	file := c.Param("file")

	p, err := checkAndJoinPath(dbName, dir, file)
	if err != nil {
		return c.String(ceIO, err.Error())
	}

	if !checkPermission(c, "api.Delete", dbName) {
		return permissionDenied(c)
	}

	err = os.Remove(p)
	if err != nil {
		errStr := fmt.Sprintf("[api.Delete] os.Remove(): %s", err.Error())
		log.Err(errStr)
		return c.String(ceIO, errStr)
	}

	return c.NoContent(200)
}

func ReadDir(c echo.Context) error {
	dbName := c.Param("db")
	dir := c.Param("dir")

	p, err := checkAndJoinPath(dbName, dir)
	if err != nil {
		return c.String(cePath, err.Error())
	}

	if !checkPermission(c, "api.ReadDir", dbName) {
		return permissionDenied(c)
	}

	files, err := os.ReadDir(p)
	if err != nil {
		errStr := fmt.Sprintf("[api.ReadDir] os.ReadDir(): %s", err.Error())
		log.Err(errStr)
		return c.String(ceIO, errStr)
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

	p, err := checkAndJoinPath(dbName)
	if err != nil {
		return c.String(cePath, err.Error())
	}

	if !checkPermission(c, "api.ReadDB", dbName) {
		return permissionDenied(c)
	}

	dirs, err := os.ReadDir(p)
	if err != nil {
		errStr := fmt.Sprintf("[api.ReadDB] os.ReadDir(): %s", err.Error())
		log.Err(errStr)
		return c.String(ceIO, errStr)
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

	p, err := checkAndJoinPath(dbName)
	if err != nil {
		return c.String(cePath, err.Error())
	}

	if !checkPermission(c, "api.DeleteDB", dbName) {
		return permissionDenied(c)
	}

	err = os.RemoveAll(p)
	if err != nil {
		errStr := fmt.Sprintf("[api.DeleteDB] os.RemoveAll(): %s", err.Error())
		log.Err(errStr)
		return c.String(ceIO, errStr)
	}

	return c.NoContent(200)
}

func DeleteDir(c echo.Context) error {
	dbName := c.Param("db")
	dir := c.Param("dir")

	p, err := checkAndJoinPath(dbName, dir)
	if err != nil {
		return c.String(cePath, err.Error())
	}

	if !checkPermission(c, "api.DeleteDir", dbName) {
		return permissionDenied(c)
	}

	err = os.RemoveAll(p)
	if err != nil {
		errStr := fmt.Sprintf("[api.DeleteDir] os.RemoveAll(): %s", err.Error())
		log.Err(errStr)
		return c.String(ceIO, errStr)
	}

	return c.NoContent(200)
}
