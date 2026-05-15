package server

import (
	"net/http"
	middleware "notes-app/internal/middleware"
	"strings"
)

type userRoutesHandler interface {
	Create(http.ResponseWriter, *http.Request)
	Show(http.ResponseWriter, *http.Request, string)
	Update(http.ResponseWriter, *http.Request, string)
}

type authRoutesHandler interface {
	Login(http.ResponseWriter, *http.Request)
}

type meRoutesHandler interface {
	Me(http.ResponseWriter, *http.Request)
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

func registerAuthRoutes(mux *http.ServeMux, authHandler authRoutesHandler) {
	mux.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/auth/login" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		authHandler.Login(w, r)
	})
}

func registerMeRoutes(
	mux *http.ServeMux,
	meHandler meRoutesHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	mux.Handle(
		"/auth/me",
		authMiddleware.RequireAuth(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
					return
				}
				meHandler.Me(w, r)
			}),
		),
	)
}
