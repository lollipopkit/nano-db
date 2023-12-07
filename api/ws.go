package api

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/lollipopkit/gommon/log"
	"github.com/lollipopkit/nano-db/cfg"
	"github.com/lollipopkit/nano-db/cst"
)

var (
	upgrader       = websocket.Upgrader{}
	errOnlyBinMsg  = errors.New("only binary message is allowed")
	errEmptyWsInfo = errors.New("empty ws info")
)

type wsInfo struct {
	Db   string
	Dir  string
	File string
	Type wsInfoType
}
type wsInfoType uint
const (
	wsInfoTypeRead wsInfoType = iota
	wsInfoTypeWrite
	wsInfoTypeDelete
)

func WS(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			return err
		}

		switch msgType {
		case websocket.TextMessage:
			return errOnlyBinMsg
		case websocket.BinaryMessage:
			// 第一行是请求信息，第二行开始是数据
			// 读取请求信息
			var info wsInfo
			var firstLine []byte
			for idx, b := range msg {
				if b == 0 {
					firstLine = msg[:idx]
					msg = msg[idx+1:]
					break
				}
			}
			if len(firstLine) == 0 {
				sendErr(errEmptyWsInfo, ws)
				continue
			}

			err = json.Unmarshal(firstLine, &info)
			if err != nil {
				sendErr(err, ws)
				continue
			}

			var paths []string
			if len(info.Db) > 0 {
				paths = append(paths, info.Db)
				if len(info.Dir) > 0 {
					paths = append(paths, info.Dir)
					if len(info.File) > 0 {
						paths = append(paths, info.File)
					}
				}
			} else {
				sendErr(errEmptyWsInfo, ws)
				continue
			}

			sn := c.Request().Header.Get(cst.HeaderKey)
			if len(sn) != cfg.App.Security.TokenLen {
				return permissionDenied(c)
			}

			if !cfg.Acl.Can(info.Db, sn) {
				return permissionDenied(c)
			}

			// 检查请求信息
			p, err := checkAndJoinPath(paths...)
			if err != nil {
				sendErr(err, ws)
				continue
			}

			switch info.Type {
			case wsInfoTypeRead:
				if info.File == "" {
					if info.Dir == "" {
						// 读取数据库
						data, err := os.ReadFile(p)
						if err != nil {
							sendErr(err, ws)
							continue
						}
						err = ws.WriteMessage(websocket.BinaryMessage, data)
						if err != nil {
							sendErr(err, ws)
							continue
						}
						continue
					}
					// 读取目录
					files, err := os.ReadDir(p)
					if err != nil {
						sendErr(err, ws)
						continue
					}
					data, err := json.Marshal(files)
					if err != nil {
						sendErr(err, ws)
						continue
					}
					err = ws.WriteMessage(websocket.BinaryMessage, data)
					if err != nil {
						sendErr(err, ws)
						continue
					}
					continue
				}
				// 读取文件
				data, err := os.ReadFile(p)
				if err != nil {
					sendErr(err, ws)
					continue
				}
				err = ws.WriteMessage(websocket.BinaryMessage, data)
				if err != nil {
					sendErr(err, ws)
					continue
				}
			case wsInfoTypeWrite:
				if info.File == "" {
					sendErr(errors.New("cannot write to db or dir"), ws)
					continue
				}
				// 写入文件
				err = os.WriteFile(p, msg, cst.FilePermission)
				if err != nil {
					sendErr(err, ws)
					continue
				}
			case wsInfoTypeDelete:
				// 删除
				err = os.RemoveAll(p)
				if err != nil {
					sendErr(err, ws)
					continue
				}
			}

		case websocket.CloseMessage:
			return nil
		case websocket.PongMessage:
			err = ws.WriteMessage(websocket.PongMessage, msg)
			if err != nil {
				return err
			}
		case websocket.PingMessage:
			err = ws.WriteMessage(websocket.PingMessage, msg)
			if err != nil {
				return err
			}
		}
	}
}

func sendErr(err error, ws *websocket.Conn) {
	log.Err(err.Error())
	ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
}
