package config

import (
	"path/filepath"
	"runtime"

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

// Config struct base from structure's config file
type Config struct {
	Server   HTTP `mapstructure:"http"`
	Iframely API  `mapstructure:"api"`
	DB       DB   `mapstructure:"db"`
}

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

// Load config from input file (yaml)
func Load() (c Config, err error) {
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(basepath)

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&c)
	return
}
