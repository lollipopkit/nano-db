package utils

import (
	"math/rand"
	"os"
	"strings"

	"git.lolli.tech/lollipopkit/nano-db/consts"
)

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func InitSalt() {
	if IsExist(consts.SaltFile) {
		salt, err := os.ReadFile(consts.SaltFile)
		if err != nil {
			println("[initSalt] os.ReadFile(): " + err.Error())
			println("[initSalt] will use default salt")
			return
		}
		consts.CookieSalt = strings.Trim(string(salt), "\n")
		return
	}
	runes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	salt := make([]rune, consts.SaltDefaultLen)
	for i := 0; i < consts.SaltDefaultLen; i++ {
		salt[i] = runes[rand.Intn(len(runes))]
	}
	os.WriteFile(consts.SaltFile, []byte(string(salt)), consts.FilePermission)
	consts.CookieSalt = string(salt)
}
