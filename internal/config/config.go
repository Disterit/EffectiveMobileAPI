package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	Database   `yaml:"database"`
	HttpServer `yaml:"HttpServer"`
}

type Database struct {
	Host     string `yaml:"host" env-default:"db"`
	user     string `yaml:"user" env-default:"postgres"`
	password string `yaml:"password" env-default:"postgres"`
	port     string `yaml:"port" env-default:"5432"`
	dbname   string `yaml:"dbname" env-default:"EffectiveMobileAPI"`
}

type HttpServer struct {
	Port        string `yaml:"port" env-default:"8080"`
	Timeout     int    `yaml:"timeout" env-default:"4"`
	IdleTimeout int    `yaml:"idleTimeout" env-default:"60"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("CONFIG_PATH does not exist")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("Error reading config file", err)
	}

	return &cfg
}
