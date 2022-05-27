package api

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/LollipopKit/nano-db/consts"
	"github.com/labstack/echo"
)

func resp(c echo.Context, code int, body interface{}) error {
	return c.JSON(200, map[string]interface{}{
		"code": code,
		"body": body,
	})
}

func path(db, col, id string) string {
	return fmt.Sprintf(pathFmt, db, col, id)
}

func accountVerify(c echo.Context) (bool, string) {
	cookieSign, errSign := c.Cookie(consts.CookieSignKey)
	cookieName, errName := c.Cookie(consts.CookieNameKey)
	if errSign != nil || errName != nil {
		return false, consts.HackUser
	}
	userName := decodeBase64(cookieName.Value)
	if userName == consts.AnonymousUser {
		return false, consts.AnonymousUser
	}
	if cookieSign.Value == generateCookieMd5(userName) {
		return true, userName
	}
	return false, consts.HackUser
}

func encodeBase64(name string) string {
	return base64.StdEncoding.EncodeToString([]byte(name))
}

func decodeBase64(b64 string) string {
	b, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		log.Println(err)
	}
	return string(b)
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
	return fmt.Sprintf("n=%s; s=%s\n", encodeBase64(userName), generateCookieMd5(userName))
}

func verifyId(id string) bool {
	if strings.Contains(id, "..") {
		return false
	}
	if strings.Contains(id, "/") {
		return false
	}
	if strings.Contains(id, "\\") {
		return false
	}
	if strings.Contains(id, " ") {
		return false
	}
	if len(id) > consts.MaxIdLength {
		return false
	}
	return true
}

func verifyParams(params []string) bool {
	for _, param := range params {
		if !verifyId(param) {
			return false
		}
	}
	return true
}