package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Docker   Docker
	GRPC     GRPC
	Logging  Logging
	Postgres Postgres
	Storage  Storage
	SSO      SSO
}

func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv() // только для локальной разработки
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func MustLoad(path string) *Config {
	cfg, err := Load(path)
	if err != nil {
		panic(err)
	}
	return cfg
}
