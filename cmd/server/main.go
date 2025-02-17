package main

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
    "shifa/internal/api"
    "shifa/pkg/database"
    "github.com/sirupsen/logrus"
    "shifa/internal/api/middleware"
)

func main() {
    // Initialize logger
    log := logrus.New()
    log.SetLevel(logrus.InfoLevel)
    log.SetFormatter(&logrus.JSONFormatter{
        TimestampFormat: "2006-01-02 15:04:05",
        PrettyPrint:    true,
    })

    // Read environment variables
    dbHost := os.Getenv("DB_HOST")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

    // Construct the connection string
    databaseURL := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbName)

    // JWT Secret - consider moving this to environment variable
    jwtSecret := "your-secret-key-here"

    // Connect to database
    db, err := database.NewMySQLConnection(databaseURL)
    if err != nil {
        log.WithError(err).Fatal("Failed to connect to database")
    }
    defer func() {
        if err := db.Close(); err != nil {
            log.WithError(err).Error("Error closing database connection")
        }
    }()

    // Verify database connection
    if err := db.Ping(); err != nil {
        log.WithError(err).Fatal("Failed to ping database")
    }
    log.Info("Successfully connected to database")

    // Setup API routes - pass db, log, and jwtSecret
    router := api.NewRouter(db, log, jwtSecret)

    // Apply CORS middleware
    corsHandler := middleware.CORSMiddleware()(router)

    // Create HTTP server with timeouts
    srv := &http.Server{
        Addr:         fmt.Sprintf(":%d", 8888),
        Handler:      corsHandler,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    // Start server in a goroutine
    go func() {
        log.WithFields(logrus.Fields{
            "port": 8888,
            "addr": srv.Addr,
        }).Info("Starting server")

        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.WithError(err).Fatal("Failed to start server")
        }
    }()

    // Setup graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    // Wait for interrupt signal
    sig := <-quit
    log.WithField("signal", sig.String()).Info("Shutting down server...")

    // Create shutdown context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Attempt graceful shutdown
    if err := srv.Shutdown(ctx); err != nil {
        log.WithError(err).Error("Server forced to shutdown")
        os.Exit(1)
    }

    log.Info("Server gracefully stopped")
}