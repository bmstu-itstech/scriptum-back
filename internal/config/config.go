package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Docker   Docker   `mapstructure:"docker"`
	GRPC     GRPC     `mapstructure:"grpc"`
	Logging  Logging  `mapstructure:"logging"`
	Postgres Postgres `mapstructure:"postgres"`
	Storage  Storage  `mapstructure:"storage"`
}

type Docker struct {
	ImagePrefix   string        `mapstructure:"image_prefix"`
	RunnerTimeout time.Duration `mapstructure:"runner_timeout"`
}

type GRPC struct {
	Port int `mapstructure:"port"`
}

type Logging struct {
	Level string `mapstructure:"level"`
}

type Postgres struct {
	URI string `mapstructure:"uri"`
}

type Storage struct {
	BasePath string `mapstructure:"base_path"`
}

func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetEnvPrefix("SC")
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
