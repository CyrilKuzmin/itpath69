// Package config application config
package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
)

// Config struct contains application configuration parameters
type Config struct {
	Host   string `env:"HOST"`
	Port   string `env:"PORT,default=1323"`
	Env    string `env:"ENV,default=local"`
	Secret string `env:"SECRET,default=secret"`
	Mongo  struct {
		URI      string `env:"URI,default=mongodb://localhost:27017/"`
		Database string `env:"SECRET,default=itpath69"`
	} `env:",prefix=MONGO_"`
}

var c *Config

// Get parses environment variables and returns fulfilled Config struct
func Get() *Config {
	if c != nil {
		return c
	}
	c = &Config{}
	if err := envconfig.Process(context.Background(), c); err != nil {
		zap.L().Fatal("parse config", zap.Error(err))
	}
	return c
}
