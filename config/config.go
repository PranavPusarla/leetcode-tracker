package config

import (
	"log/slog"
	"os"
	"time"
)

import "gopkg.in/yaml.v3"


type databaseConfig struct {
	Host string `yaml:"host"`
	Port int `yaml:"port"`
	Username string `yaml:"username"`
	Password *string `yaml:"password"`
	Name string `yaml:"name"`
}

type Config struct {
	StartDate time.Time `yaml:"start_date"`
	Users map[string]string `yaml:"users"`
	Database databaseConfig `yaml:"database"`
}

func Load(filepath string) Config {
	var data, readErr = os.ReadFile(filepath)
	if readErr != nil {
		slog.Error("Failed to read config file", "filepath", filepath)
		panic(readErr)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		slog.Error("Failed to deserialize YAML from config file", "filepath", filepath)
		panic(err)
	}
	return config
}
