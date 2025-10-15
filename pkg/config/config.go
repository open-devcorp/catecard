package config

import "os"

type Config struct {
	DatabasePath string
	ServerPort   string
	LogLevel     string
}

func LoadConfig() *Config {
	return &Config{
		// Default DB path for this project
		DatabasePath: getEnv("DB_PATH", "./catecard.db"),
		ServerPort:   getEnv("PORT", "3000"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
