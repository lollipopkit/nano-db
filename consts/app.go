package consts

const (
	FilePermission = 0770

	DBDir     = ".db/"
	LogDir    = ".log/"
	SecretDir = ".sct/"
	ACLFile   = SecretDir + "acl.json"
	SaltFile  = SecretDir + "salt.txt"

	CookieSalt2   = "20001110"
	CookieSignKey = "s"
	CookieNameKey = "n"

	AnonymousUser = "anonymous"
	HackUser      = "hack"

	MaxIdLength    = 37
	SaltDefaultLen = 17

	MaxIPFailedTimes = 37
)

var (
	// adjust this value according to your memory size.
	// bigger for better performance.
	CacherMaxLength = 100
	CookieSalt      = "nano-db"
)
