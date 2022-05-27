package model

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"git.lolli.tech/LollipopKit/nano-db/consts"
	"git.lolli.tech/LollipopKit/nano-db/utils"
)

type ACL struct {
	Version int `json:"ver"`
	Rules   []ACLRule	`json:"rules"`
}

type ACLRule struct {
	DBName string 	`json:"db"`
	UserName string	`json:"user"`
}

func (acl *ACL) Save() error {
	data, err := json.Marshal(acl)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(consts.ACLFile, data, consts.FilePermission)
}

func (acl *ACL) Load() error {
	if !utils.IsExist(consts.ACLFile) {
		err := os.MkdirAll(consts.ACLDir, consts.FilePermission)
		if err != nil {
			return err
		}

		acl = &ACL{
			Version: 1,
			Rules: []ACLRule{},
		}

		return acl.Save()
	}
	
	data, err := ioutil.ReadFile(consts.ACLFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, acl)
}

func (acl *ACL) UpdateRule(dbName, userName string) error {
	for _, rule := range acl.Rules {
		if rule.DBName == dbName {
			return nil
		}
	}
	acl.Rules = append(acl.Rules, ACLRule{
		DBName: dbName,
		UserName: userName,
	})
	return acl.Save()
}

func (acl *ACL) Can(dbName, userName string) bool {
	for _, rule := range acl.Rules {
		if rule.DBName == dbName {
			return rule.UserName == userName
		}
	}
	return false
}

func (acl *ACL) HaveDB(dbName string) bool {
	for _, rule := range acl.Rules {
		if rule.DBName == dbName {
			return true
		}
	}
	return false
}