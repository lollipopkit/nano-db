package db

import (
	"errors"
	"io/ioutil"
	"os"
	"sync"
	"time"

	glc "git.lolli.tech/lollipopkit/go_lru_cacher"
	"git.lolli.tech/lollipopkit/nano-db/consts"
	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary

	pathLockLength = consts.CacherMaxLength
	// map[string]*sync.RWMutex : {"PATH": LOCK}
	pathLockCacher = glc.NewCacher(pathLockLength)

	ErrLockConvert = errors.New("lock convert failed")
	ErrNoDocument  = errors.New("no document")
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
		// 例如：有超过pathLockLength个的进程同时读写
		for {
			if pathLockCacher.Len() < pathLockLength {
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

func Read(path string, model interface{}) error {
	return wrapLock(path, func() error {
		data, err := ioutil.ReadFile(consts.DBDir + path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return ErrNoDocument
			}
			return err
		}
		return json.Unmarshal(data, &model)
	}, false)
}

func Write(path string, model interface{}) error {
	return wrapLock(path, func() error {
		data, err := json.Marshal(model)
		if err != nil {
			return err
		}
		return ioutil.WriteFile(consts.DBDir+path, data, consts.FilePermission)
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
