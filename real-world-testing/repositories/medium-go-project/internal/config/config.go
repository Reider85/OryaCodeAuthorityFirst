package config

import (
	"os"
	"time"
)

type Config struct {
	Server   ServerConfig        `yaml:"server"`
	Database DatabaseConfig      `yaml:"database"`
	Auth     AuthConfig          `yaml:"auth"`
}

type ServerConfig struct {
	Port         string        `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type AuthConfig struct {
	JWTSecret       string        `yaml:"jwt_secret"`
	JWTExpiration   time.Duration `yaml:"jwt_expiration"`
}

func Load() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Port:         "8080",
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "password",
			Database: "app",
		},
		Auth: AuthConfig{
			JWTSecret:       "your-secret-key",
			JWTExpiration:   24 * time.Hour,
		},
	}, nil
}