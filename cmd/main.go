package main

import (
	_ "igaming/docs" // This is important!
	"igaming/internal/config"
	"igaming/internal/server"
	"log"
	"net/http"
)

// Package main iGaming API
// @title iGaming API
// @version 1.0
// @description Gaming tournament management system
// @host localhost:8080
// @BasePath /
func main() {
    cfg := config.LoadConfig()

    db := config.InitDB(cfg)
    defer db.Close()

    router := server.NewRouter(db)

    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}