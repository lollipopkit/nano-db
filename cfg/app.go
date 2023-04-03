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
			Format: "${time_rfc3339} ${method} ${status} ${uri} \nLatency: ${latency_human}  ${error}\n",
			SkipRegExp: []string{
				"/static",
				"/favicon",
			},
		},
		Security: SecurityConfig{
			Salt:      "",
			RateLimit: rate.Limit(20),
			CORSList: []string{},
			BodyLimit: "1M",
		},
	}
)

func init() {
	err := Cfg.Load()
	if err != nil {
		panic("AppConfig.Load(): " + err.Error())
	}
	if Cfg.Cache.ActiveRate < 0 || Cfg.Cache.ActiveRate > 1 {
		panic("invalid cache rate")
	}

	if Cfg.Security.Salt != "" {
		return
	}
	runes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	salt := make([]rune, consts.SaltDefaultLen)
	for i := 0; i < consts.SaltDefaultLen; i++ {
		salt[i] = runes[rand.Intn(len(runes))]
	}
	Cfg.Security.Salt = string(salt)
	err = Cfg.Save()
	if err != nil {
		panic("AppConfig.Save(): " + err.Error())
	}
}

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
	// 0.77 -> 77% of cache will be active
	// Must 0 < ActiveRate <= 1
	ActiveRate float64 `json:"active_rate"`
}

type SecurityConfig struct {
	Salt string `json:"salt"`
	// 20 -> 20 req / sec
	RateLimit rate.Limit `json:"rate_limit"`
	CORSList []string   `json:"cors_list"`
	// Limit can be specified as `4x` or `4xB`, where x is one of the multiple from K, M, G, T or P.
	BodyLimit string   `json:"body_limit"`
}

// false for better performance
// Only for echo web framework
type LogConfig struct {
	Enable bool `json:"enable"`
	// More info please refer to `echo web`
	Format string `json:"format"`
	SkipRegExp []string `json:"skip_regexp_list"`
}

func (c *AppConfig) Save() error {
	bytes, err := Json.Marshal(c)
	if err != nil {
		return errors.Join(ErrConfig, err)
	}
	if err = os.MkdirAll(consts.CfgDir, consts.FilePermission); err != nil {
		return errors.Join(ErrConfig, err)
	}
	return os.WriteFile(consts.AppCfgFile, bytes, consts.FilePermission)
}

func (c *AppConfig) Load() error {
	if !util.Exist(consts.AppCfgFile) {
		return c.Save()
	}
	bytes, err := os.ReadFile(consts.AppCfgFile)
	if err != nil {
		return errors.Join(ErrConfig, err)
	}
	return Json.Unmarshal(bytes, c)
}
