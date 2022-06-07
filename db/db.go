package db

import (
	"errors"
	"io/ioutil"
	"os"
	"sync"

	"git.lolli.tech/lollipopkit/nano-db/consts"
	"git.lolli.tech/lollipopkit/nano-db/model"
	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	// pathLocks : {"PATH": LOCK}
	pathLockCacher = model.NewCacher(consts.CacherMaxLength)
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
		lock := sync.RWMutex{}
		pathLockCacher.Update(path, &lock)
		return &lock, nil
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
			return ErrNoDocument
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
	}, false)
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
	return err
}
