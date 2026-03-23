package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Schema struct {
	Environment string `env:"ENV"`
	HttpPort    int    `env:"PORT"`
	DatabaseURI string `env:"DATABASE_URL"`
	AuthSecret  string `env:"AUTH_SECRET"`
}

const (
	ProductionEnv = "production"

	DatabaseTimeout    = 5 * time.Second
	ProductCachingTime = 1 * time.Minute
)

var (
	cfg Schema
)

func LoadConfig() *Schema {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error on load configuration file, error: %v", err)
	}
	cfg.DatabaseURI = os.Getenv("DATABASE_URL")
	if cfg.DatabaseURI == "" {
		log.Print("warning: DATABASE_URL is empty")
	}

	portStr := os.Getenv("PORT")
	if portStr == "" {
		cfg.HttpPort = 8080
	} else {
		if p, err := strconv.Atoi(portStr); err != nil {
			log.Printf("invalid PORT %q, using default 8080: %v", portStr, err)
			cfg.HttpPort = 8080
		} else {
			cfg.HttpPort = p
		}
	}

	cfg.Environment = os.Getenv("ENV")
	if cfg.Environment == "" {
		cfg.Environment = "development"
	}

	return &cfg
}

func GetEnv() *Schema {
	return &cfg
}
