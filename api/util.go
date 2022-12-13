package api

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"git.lolli.tech/lollipopkit/nano-db/consts"
	"git.lolli.tech/lollipopkit/nano-db/logger"
	"github.com/labstack/echo/v4"
)

func resp(c echo.Context, code int, data any) error {
	return c.JSON(200, map[string]any{
		"data": data,
		"code": code,
	})
}

func ok(c echo.Context) error {
	return resp(c, 200, "ok")
}

func permissionDenied(c echo.Context) error {
	return resp(c, 403, "permission denied")
}

func path(db, col, id string) string {
	return fmt.Sprintf(pathFmt, db, col, id)
}

func accountVerify(c echo.Context) (bool, string) {
	cookieSign, errSign := c.Cookie(consts.CookieSignKey)
	cookieName, errName := c.Cookie(consts.CookieNameKey)
	if errSign != nil || errName != nil {
		return false, consts.AnonymousUser
	}
	userName, err := decodeBase64(cookieName.Value)
	if err == nil && cookieSign.Value == generateCookieMd5(userName) {
		return true, userName
	}
	logger.W("[accountVerify] new hack user [%s]", userName)
	return false, consts.HackUser
}

func encodeBase64(name string) string {
	return base64.StdEncoding.EncodeToString([]byte(name))
}

func decodeBase64(b64 string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(b64)
	return string(b), err
}

func hex2Str(b [16]byte) string {
	return hex.EncodeToString(b[:])
}

func getMd5(b []byte) [16]byte {
	return md5.Sum(b)
}

func diyEncrypt(pass string) string {
	return pass[11:13] + pass + pass[4:7]
}

func generateCookieMd5(name string) string {
	base64Str := encodeBase64(consts.CookieSalt + name + consts.CookieSalt2)
	md5Hex := getMd5([]byte(base64Str))
	return diyEncrypt(reverseString(hex2Str(md5Hex)))
}

func reverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

func GenCookie(userName string) string {
	if userName == consts.AnonymousUser || userName == consts.HackUser {
		return "username cant be \"" + userName + "\", its a inner user"
	}
	return fmt.Sprintf("n=%s; s=%s", encodeBase64(userName), generateCookieMd5(userName))
}

func verifyId(id string) error {
	if strings.Contains(id, "..") {
		return errors.New(id + " contains ..")
	}
	if strings.Contains(id, "/") {
		return errors.New(id + " contains /")
	}
	if strings.Contains(id, "\\") {
		return errors.New(id + " contains \\")
	}
	if strings.Contains(id, " ") {
		return errors.New(id + " contains space")
	}
	if len(id) > consts.MaxIdLength {
		return errors.New(id + " is too long")
	}
	return nil
}

func verifyParams(params []string) error {
	for _, param := range params {
		if err := verifyId(param); err != nil {
			return err
		}
	}
	return nil
}

func checkPermission(c echo.Context, action, dbName, path string) bool {
	loggedIn, userName := accountVerify(c)

	if !loggedIn {
		return false
	}

	AclLock.RLock()
	if !Acl.Can(dbName, userName) {
		logger.W("[%s] user [%s] is trying access [%s]\n", action, userName, path)
		AclLock.RUnlock()
		return false
	}
	AclLock.RUnlock()
	return true
}
