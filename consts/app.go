package consts

const (
	// adjust this value according to your memory size.
	// bigger for better performance.
	CacherMaxLength = 100
	FilePermission  = 0770

	DBDir  = ".db/"
	LogDir = ".log/"
	ACLDir = ".acl/"
	ACLFile = ACLDir + "acl.json"

	CookieSalt    = "nano-db"
	CookieSalt2   = "20001110"
	CookieSignKey = "s"
	CookieNameKey = "n"
	CookieNotChanged = "\n!!!Attention!!!\nIt's highly recommended to change the cookie salt to your own fixed string.\n"

	AnonymousUser = "anonymous"
	HackUser      = "hack"

	MaxIdLength = 37
)
