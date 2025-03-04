package db_configs

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const pathToDBs string = "/Users/rudadadadada/diploma/db_configs"

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBname   string `yaml:"dbname"`
	Sslmode  string `yaml:"sslmode" env-default:"disable"`
}

func (cfg *DBConfig) readConfig(configPath string) *DBConfig {
	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	
	return cfg
}

type DBConfigs struct {
    Admin *DBConfig
	Client *DBConfig
}

func (cfgs *DBConfigs) GetAllDBsConfigs() *DBConfigs {
	var services = []string{"admin.yaml", "client.yaml"}

	for _, service := range services {
		var cfg DBConfig

		configPath := pathToDBs + "/" + service
		cfg.readConfig(configPath)

		switch {
		case service == "admin.yaml":
			cfgs.Admin = &cfg
		case service == "client.yaml":
			cfgs.Client = &cfg
		}
    }
	
	return cfgs
}