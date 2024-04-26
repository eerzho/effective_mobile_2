package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Http     Http
	Postgres Postgres
	Jwt      Jwt
	Logger   Logger
}

type Http struct {
	Address     string        `env:"HTTP_ADDRESS"`
	Timeout     time.Duration `env:"HTTP_TIMEOUT"`
	IdleTimeout time.Duration `env:"HTTP_IDLE_TIMEOUT"`
}

type Postgres struct {
	Url string `env:"POSTGRES_URL"`
}

type Jwt struct {
	SecretKey string `env:"JWT_SECRET_KEY"`
}

type Logger struct {
	Level string `env:"LOG_LEVEL"`
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
