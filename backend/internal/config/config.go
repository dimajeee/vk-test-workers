package config

import (
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ServerPort     string `yaml:"server_port"`
	InitialWorkers int    `yaml:"initial_workers"`
	QueueSize      int    `yaml:"queue_size"`
	LogLevel       string `yaml:"log_level"`
}

func MustLoad(path string) Config {
	cfg := Config{
		ServerPort:     "8080",
		InitialWorkers: 5,
		QueueSize:      10000,
		LogLevel:       "INFO",
	}

	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("cannot open config file: %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("cannot read config file: %v", err)
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("cannot parse config file: %v", err)
	}

	return cfg
}
