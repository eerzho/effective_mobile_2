package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Http     Http
	Postgres Postgres
	Logger   Logger
	Api      Api
}

type Http struct {
	Address string `env:"HTTP_ADDRESS"`
}

type Postgres struct {
	Url string `env:"POSTGRES_URL"`
}

type Logger struct {
	Level string `env:"LOG_LEVEL"`
}

type Api struct {
	CarInfo string `env:"API_CAR_INFO"`
}

var cfg Config

func Parse() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return err
	}

	return nil
}

func Cfg() *Config {
	return &cfg
}
