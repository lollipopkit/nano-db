package cfg

import (
	"os"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/lollipopkit/gommon/log"
	"github.com/lollipopkit/gommon/sys"
	"github.com/lollipopkit/nano-db/cst"
)

var (
	Acl = &ACL{
		Version: 1,
		Rules:   []ACLRule{},
	}
	aclLock = &sync.RWMutex{}

	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

func init() {
	go func() {
		for {
			err := Acl.load()
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
	Token   string   `json:"token"`
	DBNames []string `json:"dbs"`
}

func (acl *ACL) Save() error {
	data, err := json.MarshalIndent(acl, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cst.AclCfgFile, data, cst.FilePermission)
}

func (acl *ACL) load() error {
	aclLock.Lock()
	defer aclLock.Unlock()
	if !sys.Exist(cst.AclCfgFile) {
		err := os.MkdirAll(cst.CfgDir, cst.FilePermission)
		if err != nil {
			return err
		}

		return acl.Save()
	}

	data, err := os.ReadFile(cst.AclCfgFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, acl)
}

func (acl *ACL) updateRule(dbName, token string) error {
	aclLock.Lock()
	defer aclLock.Unlock()
	for i, rule := range acl.Rules {
		if rule.Token == token {
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
		DBNames: []string{dbName},
		Token:   token,
	})
	return acl.Save()
}

func (acl *ACL) Can(dbName, token string) bool {
	aclLock.RLock()
	defer aclLock.RUnlock()
	for _, rule := range acl.Rules {
		if rule.Token == token {
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

func UpdateAcl(token, dbName string) {
	err := Acl.updateRule(dbName, token)
	if err != nil {
		log.Err("acl update rule: " + err.Error())
	} else {
		log.Suc("acl update rule: success")
	}
}
