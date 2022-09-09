package logger

import (
	"io"
	"log"
	"os"
	"time"

	"git.lolli.tech/lollipopkit/nano-db/consts"
)

func init() {
	if err := os.MkdirAll(consts.LogDir, consts.FilePermission); err != nil {
		panic(err)
	}
	go setup()
}

func W(format string, args ...any) {
	log.Printf("[WARN] "+format, args...)
}

func I(format string, args ...any) {
	log.Printf("[INFO] "+format, args...)
}

func E(format string, args ...any) {
	log.Printf("[ERROR] "+format, args...)
}

// Must call this func using:
// `go setup()`
func setup() {
	for {
		file := consts.LogDir + time.Now().Format("2006-01-02") + ".txt"
		logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, consts.FilePermission)
		if err != nil {
			panic(err)
		}
		multiWriter := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(multiWriter)
		time.Sleep(time.Hour)
	}
}
