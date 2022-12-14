package db

import (
	"errors"
	"os"
	"sync"
	"time"

	glc "git.lolli.tech/lollipopkit/go-lru-cacher"
	"git.lolli.tech/lollipopkit/nano-db/consts"
	. "git.lolli.tech/lollipopkit/nano-db/json"
)

var (

	// map[string]*sync.RWMutex : {"PATH": LOCK}
	pathLockCacher = glc.NewCacher(consts.CacherMaxLength)

	ErrLockConvert = errors.New("lock convert failed")
)

func init() {
	if err := os.MkdirAll(consts.DBDir, consts.FilePermission); err != nil {
		panic(err)
	}
}

func getLock(path string) (*sync.RWMutex, error) {
	l, have := pathLockCacher.Get(path)
	if !have {
		// 防止pathLockCacher因为超出最大长度，而清理可能正在使用的锁
		// 例如：有超过consts.CacherMaxLength个的进程同时读写
		for {
			if pathLockCacher.Len() < consts.CacherMaxLength {
				break
			}
			time.Sleep(time.Millisecond * 100)
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
