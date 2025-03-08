package postgres_configs

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const pathToPostgresConfigs string = "/Users/rudadadadada/diploma/postgres_configs/"

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBname   string `yaml:"dbname"`
	Sslmode  string `yaml:"sslmode" env-default:"disable"`
}

func (cfg *DBConfig) GetConfig(service string) {
	configPath := pathToPostgresConfigs + service + ".yaml"
	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
}