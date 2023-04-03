package cfg

import (
	"os"
	"sync"

	"git.lolli.tech/lollipopkit/nano-db/consts"
	. "git.lolli.tech/lollipopkit/nano-db/json"
	"github.com/lollipopkit/gommon/util"
)

var (
	Acl = &ACL{
		Version: 1,
		Rules:   []ACLRule{},
	}
	aclLock = &sync.RWMutex{}
)

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
	return os.WriteFile(consts.ACLFile, data, consts.FilePermission)
}

func (acl *ACL) Load() error {
	aclLock.Lock()
	defer aclLock.Unlock()
	if !util.Exist(consts.ACLFile) {
		err := os.MkdirAll(consts.SecretDir, consts.FilePermission)
		if err != nil {
			return err
		}

		return acl.Save()
	}

	data, err := os.ReadFile(consts.ACLFile)
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
	print("[ACL]\n  ")
	err := Acl.UpdateRule(*dbName, *userName)

	if err != nil {
		println("acl update rule: " + err.Error())
	} else {
		println("acl update rule: success")
	}
}
