// pkg/database/mysql.go

package database

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewMySQLConnection(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}