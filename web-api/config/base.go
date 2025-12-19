package config

import (
	"os"
	"sync"
	"tems-web-api/utils"
)

type Config struct {
	DBHost     string `env:"DB_HOST"`
	DBPort     int    `env:"DB_PORT"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBName     string `env:"DB_NAME"`

	ServerHost string `env:"SERVER_HOST"`
	ServerPort int    `env:"SERVER_PORT"`

	WechatAppID  string `env:"WECHART_APPID"`
	WechatSecret string `env:"WECHART_SECRET"`
}

const (
	Separator     = string(os.PathSeparator)     // 路径分隔符（分隔路径元素）
	ListSeparator = string(os.PathListSeparator) // 路径列表分隔符（分隔多个路径）
)

var (
	cfg  *Config
	once sync.Once
)

// GetGlobalConfig 获取全局配置(线程安全)
func GetGlobalConfig() *Config {
	if cfg == nil {
		panic("global config not initialized")
	}
	return cfg
}

func InitLoadGlobalConfig() (*Config, error) {
	var err error
	once.Do(func() {
		cfg = applyEnvOverrides()
	})
	return cfg, err
}

func applyEnvOverrides() *Config {
	var config = Config{
		DBHost:       "127.0.0.1",
		DBPort:       5432,
		DBUser:       "postgres",
		DBPassword:   "",
		DBName:       "temsdb",
		ServerHost:   "=0.0.0.0",
		ServerPort:   55555,
		WechatAppID:  "",
		WechatSecret: "",
	}

	// 数据库配置覆盖
	if val := os.Getenv("DB_HOST"); val != "" {
		config.DBHost = val
	}
	if val := os.Getenv("DB_PORT"); val != "" {
		config.DBPort = utils.String2Int(val)
	}
	if val := os.Getenv("DB_USER"); val != "" {
		config.DBUser = val
	}
	if val := os.Getenv("DB_PASSWD"); val != "" {
		config.DBPassword = val
	}
	if val := os.Getenv("DB_NAME"); val != "" {
		config.DBName = val
	}
	// 服务配置
	if val := os.Getenv("SERVER_HOST"); val != "" {
		config.ServerHost = val
	}
	if val := os.Getenv("SERVER_PORT"); val != "" {
		config.ServerPort = utils.String2Int(val)
	}
	// 小程序配置
	if val := os.Getenv("WECHART_APPID"); val != "" {
		config.WechatAppID = val
	}
	if val := os.Getenv("WECHART_SECRET"); val != "" {
		config.WechatSecret = val
	}
	return &config
}
