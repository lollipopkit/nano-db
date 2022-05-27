package logger

import (
	"log"
	"os"
	"time"

	"github.com/LollipopKit/nano-db/consts"
)

func init() {
	if err := os.MkdirAll(consts.LogDir, consts.FilePermission); err != nil {
		panic(err)
	}
}

func W(format string, args ...interface{}) {
	log.Printf("[WARN] "+format, args...)
}

func I(format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}

func E(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format, args...)
}

// Must call this func using:
// `go logger.Setup()`
func Setup() {
	for {
		file := consts.LogDir + time.Now().Format("2006-01-02-15") + ".txt"
		logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0744)
		if err != nil {
			panic(err)
		}
		log.SetOutput(logFile)
		time.Sleep(time.Hour)
	}
}
