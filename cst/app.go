package cst

import (
	"github.com/lollipopkit/gommon/log"
)

const (
	FilePermission = 0770

	DBDir      = ".db"
	LogDir     = ".log/"
	CfgDir     = ".cfg/"
	AclCfgFile = CfgDir + "acl.json"
	AppCfgFile = CfgDir + "app.json"

	HeaderKey = "NanoDB"

	DefaultTokenLen = 37
)

func init() {
	log.Setup(log.Config{
		FilePerm:  FilePermission,
		LogPath:   LogDir,
		PrintTime: true,
	})
}
