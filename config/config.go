package config

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/viper"
)

type HTTP struct {
	Address string `mapstructur:"address"`
	Port    string `mapstructure:"port"`
}

type API struct {
	Key []string `mapstructure:"key"`
}

type DB struct {
	Source string `mapstructure:"source"`
}

type Cache struct {
	Source string        `mapstructure:"source"`
	Age    time.Duration `mapstructure:"age"`
}
type Jaeger struct {
	Source string `mapstructure:"source"`
}

// Config struct base from structure's config file
type Config struct {
	Server       HTTP   `mapstructure:"http"`
	Iframely     API    `mapstructure:"api"`
	DB           DB     `mapstructure:"db"`
	MigrationURL string `mapstructure:"migration_url"`
	Cache        Cache  `mapstructure:"cache"`
	Jaeger       Jaeger `mapstructure:"jaeger"`
}

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

// Load config from input file (yaml)
func Load() (c Config, err error) {
	if os.Getenv("DEPLOY") == "PRODUCT" {
		viper.SetConfigName("deploy")
	} else {
		viper.SetConfigName("app")
	}
	viper.SetConfigType("yaml")
	viper.AddConfigPath(basepath)

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&c)
	return
}
