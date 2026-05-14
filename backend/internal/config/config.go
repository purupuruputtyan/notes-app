package config

import (
	"fmt"
	"os"
	"strings"
)

const defaultListenAddr = ":8080"

// Config は API プロセスが起動時に読む環境変数を表す。
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	JWTSecret  string
}

// Load は必須の環境変数を読み込み Config を組み立てる。
func Load() (*Config, error) {
	host, err := getenvRequired("DB_HOST")
	if err != nil {
		return nil, err
	}
	port, err := getenvRequired("DB_PORT")
	if err != nil {
		return nil, err
	}
	user, err := getenvRequired("DB_USER")
	if err != nil {
		return nil, err
	}
	dbName, err := getenvRequired("DB_NAME")
	if err != nil {
		return nil, err
	}
	sslMode, err := getenvRequired("DB_SSLMODE")
	if err != nil {
		return nil, err
	}
	jwtSecret, err := getenvRequired("JWT_SECRET")
	if err != nil {
		return nil, err
	}
	// パスワードはローカル trust 等で空になり得るため必須チェックはしない。
	password := os.Getenv("DB_PASSWORD")

	return &Config{
		DBHost:     host,
		DBPort:     port,
		DBUser:     user,
		DBPassword: password,
		DBName:     dbName,
		DBSSLMode:  sslMode,
		JWTSecret:  jwtSecret,
	}, nil
}

// ListenAddr は HTTP の待受アドレス。Compose 等の PORT（"8080" または ":8080"）に追従する。
func (c *Config) ListenAddr() string {
	p := strings.TrimSpace(os.Getenv("PORT"))
	if p == "" {
		return defaultListenAddr
	}
	if strings.HasPrefix(p, ":") {
		return p
	}
	return ":" + p
}

// PostgresDSN は lib/pq 向けの接続文字列を返す。
func (c *Config) PostgresDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

func getenvRequired(key string) (string, error) {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return "", fmt.Errorf("required environment variable %s is not set or empty", key)
	}
	return v, nil
}
