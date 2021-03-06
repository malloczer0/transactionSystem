package config

import (
	"encoding/json"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"log"
	"sync"
)

type Config struct {
	LogLevel         string `envconfig:"LOG_LEVEL"`
	PgURL            string `envconfig:"PG_URL"`
	PgMigrationsPath string `envconfig:"PG_MIGRATIONS_PATH"`
	HTTPAddr         string `envconfig:"HTTP_ADDR"`
	PgUser           string `envconfig:"PG_USER"`
	PgPassword       string `envconfig:"PG_PASSWORD"`
}

var (
	config Config
	once   sync.Once
)

// Get reads config from environment. Once.
func Get() *Config {
	once.Do(func() {
		err := envconfig.Process("", &config)
		if err != nil {
			log.Fatal(err)
		}
		configBytes, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Configuration:", string(configBytes))
	})
	return &config
}
