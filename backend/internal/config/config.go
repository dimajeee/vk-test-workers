package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort     string
	InitialWorkers int
	QueueSize      int
	LogLevel       string
}

// Load загружает конфигурацию из переменных окружения с дефолтами
func MustLoad() Config {
	return Config{
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		InitialWorkers: getEnvAsInt("INITIAL_WORKERS", 5),
		QueueSize:      getEnvAsInt("QUEUE_SIZE", 10000),
		LogLevel:       getEnv("LOG_LEVEL", "INFO"),
	}
}

func getEnv(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	if valStr, ok := os.LookupEnv(name); ok {
		if val, err := strconv.Atoi(valStr); err == nil {
			return val
		}
	}
	return defaultVal
}
