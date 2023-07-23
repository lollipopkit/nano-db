package cfg

import (
	"errors"
	"os"

	"github.com/lollipopkit/gommon/sys"
	"github.com/lollipopkit/nano-db/cst"
	"golang.org/x/time/rate"
)

var (
	ErrConfig = errors.New("config file error")

	App = &AppConfig{
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
			TokenLen:  17,
			RateLimit: rate.Limit(20),
			CORSList:  []string{},
			BodyLimit: "1M",
		},
		Misc: MiscConfig{
			MaxPathLen: 37,
		},
	}
)

func init() {
	err := App.Load()
	if err != nil {
		panic("AppConfig.Load(): " + err.Error())
	}
	if App.Cache.ActiveRate < 0 || App.Cache.ActiveRate > 1 {
		panic("invalid cache rate")
	}
	if App.Misc.MaxPathLen <= 0 {
		panic("invalid max path len")
	}
	if App.Security.TokenLen <= 0 {
		App.Security.TokenLen = 17
	}

	err = App.Save()
	if err != nil {
		panic("AppConfig.Save(): " + err.Error())
	}
}

type AppConfig struct {
	Addr     string         `json:"addr"`
	Cache    CacheConfig    `json:"cache"`
	Log      LogConfig      `json:"log"`
	Security SecurityConfig `json:"security"`
	Misc     MiscConfig     `json:"misc"`
}

type MiscConfig struct {
	// Max len for:
	//
	//	{{DB}} {{DIR}} {{FILE}}
	MaxPathLen int `json:"max_path_len"`
}

type CacheConfig struct {
	// adjust this value according to your memory size.
	//
	// bigger for better performance.
	MaxSize int `json:"max_size"`
	// 0.77 -> 77% of cache will be active
	//
	// Must 0 < ActiveRate <= 1
	ActiveRate float64 `json:"active_rate"`
}

type SecurityConfig struct {
	// Default: 17
	TokenLen int `json:"token_len"`
	// 20 -> 20 req / sec
	RateLimit rate.Limit `json:"rate_limit"`
	CORSList  []string   `json:"cors_list"`
	// Limit can be specified as `4x` or `4xB`, where x is one of the multiple from K, M, G, T or P.
	BodyLimit string `json:"body_limit"`
}

type LogConfig struct {
	// false for better performance
	//
	// Only for echo web framework
	Enable bool `json:"enable"`
	// More info please refer to `echo web`
	Format     string   `json:"format"`
	SkipRegExp []string `json:"skip_regexp_list"`
}

func (c *AppConfig) Save() error {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return errors.Join(ErrConfig, err)
	}
	if err = os.MkdirAll(cst.CfgDir, cst.FilePermission); err != nil {
		return errors.Join(ErrConfig, err)
	}
	return os.WriteFile(cst.AppCfgFile, bytes, cst.FilePermission)
}

func (c *AppConfig) Load() error {
	if !sys.Exist(cst.AppCfgFile) {
		return c.Save()
	}
	bytes, err := os.ReadFile(cst.AppCfgFile)
	if err != nil {
		return errors.Join(ErrConfig, err)
	}
	return json.Unmarshal(bytes, c)
}
