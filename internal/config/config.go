package config

import (
	"sync"
	"vortex_test/pkg/logging"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Host        string `env:"HOST" env-default:"127.0.0.1"`
	Port        string `env:"PORT" env-default:"8080"`
	DatabaseURL string `env:"DATABASEURL" env-default:"postgres://postgres:postgres@localhost:5432/chsuBot?sslmode=disable"`
}

var (
	cfg  *Config
	once sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.New()
		logger.Info("Read app configuration")
		cfg = &Config{}
		if err := cleanenv.ReadEnv(cfg); err != nil {
			help, _ := cleanenv.GetDescription(cfg, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return cfg
}
