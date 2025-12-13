package config

import (
	"fmt"
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
	SSO      SSO      `mapstructure:"sso"`
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

type SSO struct {
	Host  string `mapstructure:"host"`
	Port  string `mapstructure:"port"`
	AppID int32  `mapstructure:"app_id"`
}

func Load(path string) (*Config, error) {
	// Нетривиальный момент Viper, не описанный в документации, но описанный в
	// 	https://github.com/spf13/viper/issues/1797
	// Без viper.ExperimentalBindStruct() переменная окружения загружается только если она была указана в yaml конфиге.
	// Так например:
	// 	postgres:
	//    uri:
	// и SC_POSTGRES_URI работает корректно, а без пустого uri -- не читает переменную вовсе. Причём,
	//	viper.Get("postgres.uri")
	// работает исправно -- проблема именно в Unmarshall, который почему-то полагается на файл.
	// Решение -- viper.ExperimentalBindStruct().
	v := viper.NewWithOptions(viper.ExperimentalBindStruct())
	v.SetConfigFile(path)
	v.AutomaticEnv()
	v.SetEnvPrefix("SC")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config '%s': %w", path, err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
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
