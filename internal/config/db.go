package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB(cfg *Config) *sql.DB {
    dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
        cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

    db, err := sql.Open("mysql", dataSource)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    //   connection test
    if err := db.Ping(); err != nil {
        log.Fatalf("Failed to ping database: %v", err)
    }

    return db
}