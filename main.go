package main

import (
	"flag"
	"math/rand"
	"os"

	"github.com/lollipopkit/gommon/log"
	. "github.com/lollipopkit/nano-db/cfg"
)

func main() {
	parseCli()
	if err := startWeb(); err != nil {
		log.Err(err.Error())
	}
}

func parseCli() {
	token := flag.String("t", "", "set token with -u <token>")
	dbName := flag.String("d", "", "update acl rules with -d <dbname>")
	flag.Parse()

	// generate cookie & update acl rules
	if *dbName != "" {
		if *token == "" {
			*token = genToken()
			log.Info("generated token: %s", *token)
		}
		UpdateAcl(*token, *dbName)
		os.Exit(0)
	}
}

func genToken() string {
	runes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	salt := make([]rune, App.Security.TokenLen)
	for i := 0; i < App.Security.TokenLen; i++ {
		salt[i] = runes[rand.Intn(len(runes))]
	}
	return string(salt)
}
