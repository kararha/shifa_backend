package main

import (
    "database/sql"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "shifa/internal/api"
    "shifa/pkg/fileutils"

    _ "github.com/go-sql-driver/mysql"
    "github.com/sirupsen/logrus"
)

func main() {
    // Initialize logger
    log := logrus.New()
    log.SetFormatter(&logrus.JSONFormatter{})

    // Connect to the database
    db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/dbname")
    if err != nil {
        log.Fatalf("Failed to connect to the database: %v", err)
    }
    defer db.Close()

    // Ensure upload directories exist
    if err := os.MkdirAll(filepath.Join(fileutils.UploadDir, fileutils.ProfileImages), 0755); err != nil {
        log.Fatalf("Failed to create upload directories: %v", err)
    }

    // Configure file upload settings
    http.MaxBytesReader = func(w http.ResponseWriter, r *http.Request, n int64) *http.Request {
        r.Body = http.MaxBytesReader(w, r.Body, n)
        return r
    }

    // Initialize router
    router := api.NewRouter(db, log, "your_jwt_secret")

    // Start the server
    log.Info("Starting server on :8080")
    if err := http.ListenAndServe(":8080", router); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}