package config

import (
	"github.com/spf13/viper"
)

type HTTP struct {
	Address string `mapstructur:"address"`
	Port    string `mapstructure:"port"`
}

type API struct {
	Key string `mapstructure:"key"`
}

// Config struct base from structure's config file
type Config struct {
	Server   HTTP `mapstructure:"http"`
	Iframely API  `mapstructure:"api"`
}

// Load config from input file (yaml)
func Load() (c Config, err error) {
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath(".")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&c)
	return
}
