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
	statusFmt              = "%d dbs, %d dirs in %s"
	contextKeyPathNotFound = "context key path not found"
)

func Alive(c echo.Context) error {
	return c.NoContent(200)
}

func Read(c echo.Context) error {
	p, ok := c.Get(contextKeyPath).(string)
	if !ok {
		return c.String(cePath, contextKeyPathNotFound)
	}

	return send(c, p)
}

func Write(c echo.Context) error {
	p, ok := c.Get(contextKeyPath).(string)
	if !ok {
		return c.String(cePath, contextKeyPathNotFound)
	}

	err := os.MkdirAll(filepath.Dir(p), cst.FilePermission)
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
	p, ok := c.Get(contextKeyPath).(string)
	if !ok {
		return c.String(cePath, contextKeyPathNotFound)
	}

	err := os.Remove(p)
	if err != nil {
		errStr := fmt.Sprintf("[api.Delete] os.Remove(): %s", err.Error())
		log.Err(errStr)
		return c.String(ceIO, errStr)
	}

	return c.NoContent(200)
}

func ReadDir(c echo.Context) error {
	p, ok := c.Get(contextKeyPath).(string)
	if !ok {
		return c.String(cePath, contextKeyPathNotFound)
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
	p, ok := c.Get(contextKeyPath).(string)
	if !ok {
		return c.String(cePath, contextKeyPathNotFound)
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
	p, ok := c.Get(contextKeyPath).(string)
	if !ok {
		return c.String(cePath, contextKeyPathNotFound)
	}

	err := os.RemoveAll(p)
	if err != nil {
		errStr := fmt.Sprintf("[api.DeleteDB] os.RemoveAll(): %s", err.Error())
		log.Err(errStr)
		return c.String(ceIO, errStr)
	}

	return c.NoContent(200)
}

func DeleteDir(c echo.Context) error {
	p, ok := c.Get(contextKeyPath).(string)
	if !ok {
		return c.String(cePath, contextKeyPathNotFound)
	}

	err := os.RemoveAll(p)
	if err != nil {
		errStr := fmt.Sprintf("[api.DeleteDir] os.RemoveAll(): %s", err.Error())
		log.Err(errStr)
		return c.String(ceIO, errStr)
	}

	return c.NoContent(200)
}
