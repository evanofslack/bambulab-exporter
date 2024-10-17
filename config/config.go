package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	App  `yaml:"app"`
	HTTP `yaml:"http"`
	Log  `yaml:"logger"`
    Auth `yaml:"auth"`
}

type App struct {
	Name    string `yaml:"name" env:"APP_NAME"`
	Version string `yaml:"version" env:"APP_VERSION"`
	Env     string `yaml:"env" env:"APP_ENV"`
}

type Auth struct {
	Username string `yaml:"username" env:"BAMBU_USERNAME"`
	Password string `yaml:"password" env:"BAMBU_PASSWORD"`
	Key      string `yaml:"apikey" env:"BAMBU_API_KEY"`
	Endpoint string `yaml:"endpoint" env:"BAMBU_ENDPOINT"`
	DeviceId string `yaml:"deviceid" env:"BAMBU_DEVICE_ID"`
}

type HTTP struct {
	Port string `yaml:"port" env:"HTTP_PORT"`
}

type Log struct {
	Level string `yaml:"level" env:"LOG_LEVEL"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := godotenv.Load(); err != nil {
		fmt.Println("Could not load .env file")
	}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("Error loading env: %w", err)
	}
	return cfg, nil
}
