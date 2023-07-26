package config

import (
	"L0/internal/database"
	"github.com/spf13/viper"
)

type Env struct {
	Config database.Config
}

func Init() (*Env, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Env

	if err := viper.UnmarshalKey("pg_config", &cfg.Config); err != nil {
		return nil, err
	}
	return &cfg, nil
}
