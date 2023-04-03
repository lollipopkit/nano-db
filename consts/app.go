package consts

import "github.com/lollipopkit/gommon/term"

const (
	FilePermission = 0770

	DBDir      = ".db/"
	LogDir     = ".log/"
	CfgDir     = ".cfg/"
	AclCfgFile = CfgDir + "acl.json"
	AppCfgFile = CfgDir + "app.json"

	CookieSalt2   = "20001110"
	CookieSignKey = "s"
	CookieNameKey = "n"

	MaxIdLength    = 37
	SaltDefaultLen = 17
)

func init() {
	term.SetLog(LogDir, FilePermission, true)
}
