package config

import (
    "fmt"
    "os"
    "strconv"
)

type Config struct {
    DatabaseURL string
    ServerPort  int
    LogLevel    string
    LogFormat   string
}

func Load() (*Config, error) {
    fmt.Println("=== Loading Configuration ===")
    
    // Load environment variables
    dbHost := os.Getenv("DB_HOST")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    serverPort := os.Getenv("SERVER_PORT")
    logLevel := os.Getenv("LOG_LEVEL")
    logFormat := os.Getenv("LOG_FORMAT")

    fmt.Printf("Environment variables:\n")
    fmt.Printf("DB_HOST: %s\n", dbHost)
    fmt.Printf("DB_USER: %s\n", dbUser)
    fmt.Printf("DB_NAME: %s\n", dbName)
    fmt.Printf("SERVER_PORT: %s\n", serverPort)

    if dbHost == "" || dbUser == "" || dbName == "" {
        return nil, fmt.Errorf("missing required environment variables for the database")
    }

    // Build the database URL
    dbURL := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", 
        dbUser, 
        dbPassword, 
        dbHost, 
        dbName,
    )

    // Set default port
    port := 8888 // default port
    if serverPort != "" {
        var err error
        port, err = strconv.Atoi(serverPort)
        if err != nil {
            return nil, fmt.Errorf("invalid server port: %v", err)
        }
    }
    fmt.Printf("Final port selected: %d\n", port)

    // Set defaults for log configuration
    if logLevel == "" {
        logLevel = "info"
    }
    if logFormat == "" {
        logFormat = "json"
    }

    config := &Config{
        DatabaseURL: dbURL,
        ServerPort:  port,
        LogLevel:    logLevel,
        LogFormat:   logFormat,
    }

    fmt.Printf("=== Final Configuration ===\n")
    fmt.Printf("ServerPort: %d\n", config.ServerPort)
    fmt.Printf("DatabaseURL: %s\n", config.DatabaseURL)
    fmt.Printf("LogLevel: %s\n", config.LogLevel)
    fmt.Printf("LogFormat: %s\n", config.LogFormat)

    return config, nil
}