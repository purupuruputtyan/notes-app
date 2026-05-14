package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	handler "notes-app/internal/handler/user"
	repository "notes-app/internal/repository/user"
	usecase "notes-app/internal/usecase/user"
)

const (
	defaultListenAddr = ":8080"
	shutdownTimeout   = 10 * time.Second
)

// listenAddr は HTTP の待受アドレス。Compose 等の PORT（"8080" または ":8080"）に追従する。
func listenAddr() string {
	p := strings.TrimSpace(os.Getenv("PORT"))
	if p == "" {
		return defaultListenAddr
	}
	if strings.HasPrefix(p, ":") {
		return p
	}
	return ":" + p
}

type userRoutesHandler interface {
	Create(http.ResponseWriter, *http.Request)
	Show(http.ResponseWriter, *http.Request, string)
	Update(http.ResponseWriter, *http.Request, string)
}

func main() {
	db, err := newDB()
	if err != nil {
		log.Fatalf("db connection failed: host=%s err=%v", os.Getenv("DB_HOST"), err)
	}
	defer db.Close()

	server := newServer(listenAddr(), db)

	go listen(server)

	waitForShutdown(server)
}

// getenvRequired は必須の環境変数を読む。未設定や空白のみのときはエラーを返す。
func getenvRequired(key string) (string, error) {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return "", fmt.Errorf("required environment variable %s is not set or empty", key)
	}
	return v, nil
}

func newDB() (*sql.DB, error) {
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
	// パスワードはローカル trust 等で空になり得るため必須チェックはしない。
	password := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbName, sslMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	// 接続プールの既定（小〜中規模 API 向け）。負荷に合わせて調整する。
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

func newServer(addr string, db *sql.DB) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           newMux(db),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}

func newMux(db *sql.DB) http.Handler {
	repo := repository.NewUser(db)
	userUsecase := usecase.NewUserUseCase(repo)
	userHandler := handler.New(userUsecase)

	mux := http.NewServeMux()

	registerRootRoute(mux)
	registerHealthRoutes(mux)
	registerUserRoutes(mux, userHandler)

	return mux
}

func registerRootRoute(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"Go API is running"}`))
	})
}

func registerHealthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
}

func registerUserRoutes(mux *http.ServeMux, userHandler userRoutesHandler) {
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			userHandler.Create(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/users/")
		if id == "" {
			http.NotFound(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			userHandler.Show(w, r, id)
		case http.MethodPut:
			userHandler.Update(w, r, id)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func listen(server *http.Server) {
	log.Printf("API server started on %s", server.Addr)

	if err := server.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func waitForShutdown(server *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		shutdownTimeout,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}

	log.Println("server stopped")
}
