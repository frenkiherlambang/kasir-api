package config

import (
	"errors"

	"github.com/spf13/viper"
)

// Config holds application configuration.
type Config struct {
	DBConn string
	Port   string
}

// Load reads configuration from .env and environment variables.
// Environment variables override values from the file.
// Returns an error if DB_CONN is empty.
func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	_ = viper.ReadInConfig() // ignore file-not-found; env vars still work

	cfg := &Config{
		DBConn: viper.GetString("DB_CONN"),
		Port:   viper.GetString("PORT"),
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.DBConn == "" {
		return nil, errors.New("DB_CONN is required (set in .env or environment)")
	}
	return cfg, nil
}
