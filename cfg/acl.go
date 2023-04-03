package cfg

import (
	"os"
	"sync"
	"time"

	"github.com/lollipopkit/gommon/term"
	"github.com/lollipopkit/gommon/util"
	"github.com/lollipopkit/nano-db/consts"
	. "github.com/lollipopkit/nano-db/json"
)

var (
	Acl = &ACL{
		Version: 1,
		Rules:   []ACLRule{},
	}
	aclLock = &sync.RWMutex{}
)

func init() {
	go func() {
		for {
			err := Acl.Load()
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Minute)
		}
	}()
}

type ACL struct {
	Version int       `json:"ver"`
	Rules   []ACLRule `json:"rules"`
}

type ACLRule struct {
	UserName string   `json:"user"`
	DBNames  []string `json:"db"`
}

func (acl *ACL) Save() error {
	data, err := Json.Marshal(acl)
	if err != nil {
		return err
	}
	return os.WriteFile(consts.AclCfgFile, data, consts.FilePermission)
}

func (acl *ACL) Load() error {
	aclLock.Lock()
	defer aclLock.Unlock()
	if !util.Exist(consts.AclCfgFile) {
		err := os.MkdirAll(consts.CfgDir, consts.FilePermission)
		if err != nil {
			return err
		}

		return acl.Save()
	}

	data, err := os.ReadFile(consts.AclCfgFile)
	if err != nil {
		return err
	}

	return Json.Unmarshal(data, acl)
}

func (acl *ACL) UpdateRule(dbName, userName string) error {
	aclLock.Lock()
	defer aclLock.Unlock()
	for i, rule := range acl.Rules {
		if rule.UserName == userName {
			for _, db := range rule.DBNames {
				if db == dbName {
					return nil
				}
			}
			acl.Rules[i].DBNames = append(rule.DBNames, dbName)
			return acl.Save()
		}
	}
	acl.Rules = append(acl.Rules, ACLRule{
		DBNames:  []string{dbName},
		UserName: userName,
	})
	return acl.Save()
}

func (acl *ACL) Can(dbName, userName string) bool {
	aclLock.RLock()
	defer aclLock.RUnlock()
	for _, rule := range acl.Rules {
		if rule.UserName == userName {
			for _, db := range rule.DBNames {
				if db == dbName {
					return true
				}
			}
			break
		}
	}
	return false
}

func UpdateAcl(userName, dbName *string) {
	err := Acl.UpdateRule(*dbName, *userName)
	if err != nil {
		term.Err("acl update rule: " + err.Error())
	} else {
		term.Suc("acl update rule: success")
	}
}
