package consts

const (
	FilePermission = 0770

	DBDir   = ".db/"
	LogDir  = ".log/"
	CfgDir  = ".cfg/"
	ACLFile = CfgDir + "acl.json"
	CfgFile = CfgDir + "app.json"

	CookieSalt2   = "20001110"
	CookieSignKey = "s"
	CookieNameKey = "n"

	AnonymousUser = "anonymous"
	HackUser      = "hack"

	MaxIdLength    = 37
	SaltDefaultLen = 17
)
