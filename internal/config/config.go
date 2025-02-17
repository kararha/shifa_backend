package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort int
	LogLevel   logrus.Level
	LogFormat  string
	JWTSecret  string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		logrus.Warn("Error loading .env file: %v", err)
	}

	logLevelStr := os.Getenv("LOG_LEVEL")
	logLevel, err := logrus.ParseLevel(logLevelStr)
	if err != nil {
		logLevel = logrus.InfoLevel
	}

	serverPortStr := os.Getenv("SERVER_PORT")
	serverPort, err := strconv.Atoi(serverPortStr)
	if err != nil {
		serverPort = 8888
	}

	return &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		ServerPort: serverPort,
		LogLevel:   logLevel,
		LogFormat:  os.Getenv("LOG_FORMAT"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
	}, nil
}
