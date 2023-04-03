package db

import (
	"errors"
	"os"
	"sync"
	"time"

	glc "github.com/lollipopkit/go-lru-cacher"
	. "github.com/lollipopkit/nano-db/cfg"
	"github.com/lollipopkit/nano-db/consts"
	. "github.com/lollipopkit/nano-db/json"
)

var (
	// map[string]*sync.RWMutex : {"PATH": LOCK}
	pathLockCacher *glc.Cacher

	ErrLockConvert = errors.New("lock convert failed")
)

func init() {
	if err := os.MkdirAll(consts.DBDir, consts.FilePermission); err != nil {
		panic(err)
	}
	pathLockCacher = glc.NewCacher(Cfg.Cache.MaxSize)
}

func getLock(path string) (*sync.RWMutex, error) {
	l, have := pathLockCacher.Get(path)
	if !have {
		// 防止 pathLockCacher 因为超出最大长度，而清理可能正在使用的锁
		// 例如：有超过 MaxSize 个的进程同时读写
		for pathLockCacher.Len() == Cfg.Cache.MaxSize {
			time.Sleep(time.Millisecond * 17)
		}

		lock := new(sync.RWMutex)
		pathLockCacher.Set(path, lock)
		return lock, nil
	}

	a, ok := l.(*sync.RWMutex)
	if ok {
		return a, nil
	}
	return nil, ErrLockConvert
}

func Read(path string, model any) error {
	return wrapLock(path, func() error {
		data, err := os.ReadFile(consts.DBDir + path)
		if err != nil {
			return err
		}
		return Json.Unmarshal(data, &model)
	}, false)
}

func Write(path string, model any) error {
	return wrapLock(path, func() error {
		data, err := Json.Marshal(model)
		if err != nil {
			return err
		}
		return os.WriteFile(consts.DBDir+path, data, consts.FilePermission)
	}, true)
}

func Delete(path string) error {
	return wrapLock(path, func() error {
		return os.Remove(consts.DBDir + path)
	}, true)
}

func wrapLock(path string, fun func() error, write bool) error {
	lock, err := getLock(path)
	if err != nil {
		return err
	}

	if write {
		lock.Lock()
	} else {
		lock.RLock()
	}

	err = fun()

	if write {
		lock.Unlock()
	} else {
		lock.RUnlock()
	}

	pathLockCacher.Delete(path)
	return err
}
