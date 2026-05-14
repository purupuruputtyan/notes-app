package server

import (
	"database/sql"
	"net/http"

	authhandler "notes-app/internal/handler/auth"
	userhandler "notes-app/internal/handler/user"
	userrepo "notes-app/internal/repository/user"
	authusecase "notes-app/internal/usecase/auth"
	userusecase "notes-app/internal/usecase/user"
)

// NewMux は依存関係を組み立て、ルートを登録した http.Handler を返す。
func NewMux(db *sql.DB) http.Handler {
	repo := userrepo.NewUser(db)
	userUC := userusecase.NewUserUseCase(repo)
	userH := userhandler.New(userUC)
	loginUC := authusecase.NewLoginUseCase(repo)
	loginH := authhandler.NewLoginHandler(loginUC)

	mux := http.NewServeMux()

	registerRootRoute(mux)
	registerHealthRoutes(mux)
	registerUserRoutes(mux, userH)
	registerAuthRoutes(mux, loginH)

	return mux
}
