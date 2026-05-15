package main

import (
	"log"

	"notes-app/internal/config"
	"notes-app/internal/database"
	"notes-app/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := database.Open(cfg)
	if err != nil {
		log.Fatalf("db connection failed: host=%s err=%v", cfg.DBHost, err)
	}
	defer db.Close()

	mux := server.NewMux(db, cfg.JWTSecret)
	srv := server.NewHTTPServer(cfg.ListenAddr(), mux)

	server.Run(srv)
}
