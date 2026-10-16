package config

import (
	"net/url"
	"os"
)

type Config struct {
	Host       string
	Port       string
	GinMode    string
	APIPrefix  string
	APIVersion string
	DSN        string
}

func (c *Config) BasePath() string {
	return c.APIPrefix + "/" + c.APIVersion
}

func Load() *Config {
	return &Config{
		Host:       "0.0.0.0",
		Port:       "8080",
		GinMode:    "debug",
		APIPrefix:  "/api",
		APIVersion: "v1",
		DSN: (&url.URL{
			Scheme:   "postgres",
			User:     url.UserPassword(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD")),
			Host:     "postgres:5432",
			Path:     os.Getenv("POSTGRES_DB"),
			RawQuery: "sslmode=disable",
		}).String(),
	}
}
