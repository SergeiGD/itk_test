package config

import "time"

type Config struct {
	App struct {
		Name    string `yaml:"name" env:"APP_NAME"`
		Version string `yaml:"version" env:"APP_VERSION"`
		Debug   bool   `yaml:"Debug" env:"DEBUG"`
	}
	Database struct {
		Host       string        `env:"DATABASE_HOST"`
		User       string        `env:"DATABASE_USER"`
		Password   string        `env:"DATABASE_PASSWORD"`
		Port       string        `env:"DATABASE_PORT"`
		Db         string        `env:"DATABASE_NAME"`
		Timeout    time.Duration `env:"DATABASE_TIMEOUT"`
		ConnDelay  time.Duration `env:"DATABASE_DELAY"`
		MaxAttemps int           `env:"DATABASE_MAXATTEMPS"`
	}
	Limiter struct {
		MaxLimit      int           `env:"LIMITER_MAX_LIMIT"`
		Burst         int           `env:"LIMITER_BURST"`
		CleanInterval time.Duration `env:"LIMITER_CLEAN_INTERVAL"`
	}
}
