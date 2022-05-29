package utils

import "os"

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
