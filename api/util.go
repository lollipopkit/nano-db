package api

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lollipopkit/gommon/term"
	"github.com/lollipopkit/nano-db/cfg"
	"github.com/lollipopkit/nano-db/consts"
)

const (
	pathFmt = "%s/%s/%s"
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

func path(db, dir, id string) string {
	return fmt.Sprintf(pathFmt, db, dir, id)
}

func accountVerify(c echo.Context) (bool, string) {
	cookieSign, errSign := c.Cookie(consts.CookieSignKey)
	cookieName, errName := c.Cookie(consts.CookieNameKey)
	if errSign != nil || errName != nil {
		return false, ""
	}
	userName, err := decodeBase64(cookieName.Value)
	if userName == "" {
		return false, ""
	}
	if err == nil && cookieSign.Value == generateCookieMd5(userName) {
		return true, userName
	}
	term.Warn("[accountVerify] hack user [%s]", userName)
	return false, ""
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
	base64Str := encodeBase64(cfg.Cfg.Security.Salt + name + consts.CookieSalt2)
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
	if userName == "" {
		return "username cant be '" + userName + "'"
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

	if !cfg.Acl.Can(dbName, userName) {
		term.Warn("[%s] user [%s] is trying access [%s]", action, userName, path)
		return false
	}
	return true
}
