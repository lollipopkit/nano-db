package cfg

import (
	"errors"
	"math/rand"
	"os"

	"github.com/lollipopkit/gommon/util"
	"github.com/lollipopkit/nano-db/consts"
	. "github.com/lollipopkit/nano-db/json"
	"golang.org/x/time/rate"
)

var (
	ErrConfig = errors.New("config file error")

	Cfg = &AppConfig{
		Addr: ":3770",
		Cache: CacheConfig{
			MaxSize:    100,
			ActiveRate: 0.77,
		},
		Log: LogConfig{
			Enable: false,
			Debug:  false,
		},
		Security: SecurityConfig{
			Salt:      "",
			RateLimit: rate.Limit(20),
		},
	}
)

type AppConfig struct {
	Addr     string         `json:"addr"`
	Cache    CacheConfig    `json:"cache"`
	Log      LogConfig      `json:"log"`
	Security SecurityConfig `json:"security"`
}

type CacheConfig struct {
	// adjust this value according to your memory size.
	// bigger for better performance.
	MaxSize    int     `json:"max_size"`
	ActiveRate float64 `json:"active_rate"`
}

type SecurityConfig struct {
	Salt string `json:"salt"`
	// 20 -> 20 req / sec
	RateLimit rate.Limit `json:"rate_limit"`
}

// false for better performance
type LogConfig struct {
	Enable bool `json:"enable"`
	Debug  bool `json:"debug"`
}

func (c *AppConfig) Save() error {
	bytes, err := Json.Marshal(c)
	if err != nil {
		return errors.Join(ErrConfig, err)
	}
	return os.WriteFile(consts.CfgFile, bytes, consts.FilePermission)
}

func (c *AppConfig) Load() error {
	if !util.Exist(consts.CfgFile) {
		return c.Save()
	}
	bytes, err := os.ReadFile(consts.CfgFile)
	if err != nil {
		return errors.Join(ErrConfig, err)
	}
	return Json.Unmarshal(bytes, c)
}

func init() {
	if Cfg.Security.Salt != "" {
		return
	}
	runes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	salt := make([]rune, consts.SaltDefaultLen)
	for i := 0; i < consts.SaltDefaultLen; i++ {
		salt[i] = runes[rand.Intn(len(runes))]
	}
	Cfg.Security.Salt = string(salt)
}
