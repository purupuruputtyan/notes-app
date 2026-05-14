package database

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"

	"notes-app/internal/config"
)

// Open は設定に基づき Postgres に接続し、接続プールの既定を設定する。
func Open(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.PostgresDSN())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
