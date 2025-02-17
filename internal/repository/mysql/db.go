package mysql

import (
	"database/sql"
	"fmt"
	"shifa/internal/config"
	"time"

	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
	"github.com/sirupsen/logrus"
)

// NewMySQLDB creates a new MySQL database connection with retries
func NewMySQLDB(cfg *config.Config, log *logrus.Logger) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBName)

	// Retry configuration
	maxRetries := 30
	retryInterval := time.Second * 2

	var db *sql.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Warnf("Failed to open database connection (attempt %d/%d): %v", i+1, maxRetries, err)
			time.Sleep(retryInterval)
			continue
		}

		// Test the connection
		err = db.Ping()
		if err == nil {
			log.Infof("Successfully connected to database after %d attempts", i+1)
			return db, nil
		}

		log.Warnf("Failed to ping database (attempt %d/%d): %v", i+1, maxRetries, err)
		db.Close()
		time.Sleep(retryInterval)
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %v", maxRetries, err)
}
