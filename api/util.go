package api

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/lollipopkit/nano-db/cfg"
	"github.com/lollipopkit/nano-db/cst"
)

var (
	errEmptyPath   = errors.New("empty path")
	errPathTooLong = errors.New("path too long")
	errPathDot     = errors.New("path cannot start or end with '.'")
)

func checkAndJoinPath(paths ...string) (string, error) {
	for _, p := range paths {
		if err := verifyPath(p); err != nil {
			return "", fmt.Errorf("%s is invalid: %s", paths, err.Error())
		}
	}
	return filepath.Join(cst.DBDir, filepath.Join(paths...)), nil
}

// 如果包含除 0-9 A-Z a-z . 以外的字符，返回错误
func verifyPath(s string) error {
	if len(s) == 0 {
		return errEmptyPath
	}
	runes := []rune(s)
	if len(runes) > cfg.App.Misc.MaxPathLen {
		return errPathTooLong
	}
	if runes[0] == 46 || runes[len(runes)-1] == 46 {
		return errPathDot
	}
	for _, r := range runes {
		if (r >= 48 && r <= 57) ||
			(r >= 65 && r <= 90) ||
			(r >= 97 && r <= 122) ||
			r == 46 {
			continue
		}
		return fmt.Errorf("invalid character '%s'", string(r))
	}
	return nil
}
