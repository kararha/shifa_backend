package config

import (
    "os"
    "strconv"
)

type Config struct {
    JWTSecret    string
    ServerPort   int
    DatabaseURL  string
    LogLevel     string
    LogFormat    string
}

func Load() (*Config, error) {
    port, err := strconv.Atoi(getEnvOrDefault("SERVER_PORT", "8080"))
    if err != nil {
        return nil, err
    }

    return &Config{
        ServerPort:  port,
        DatabaseURL: getEnvOrDefault("DATABASE_URL", "root:@tcp(localhost:3306)/shifa?parseTime=true"),
        LogLevel:   getEnvOrDefault("LOG_LEVEL", "info"),
        LogFormat:  getEnvOrDefault("LOG_FORMAT", "json"),
    }, nil
}

func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}