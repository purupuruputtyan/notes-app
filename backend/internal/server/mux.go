package server

import (
	"database/sql"
	"net/http"

	authHandler "notes-app/internal/handler/auth"
	userHandler "notes-app/internal/handler/user"
	userRepo "notes-app/internal/repository/user"
	authUsecase "notes-app/internal/usecase/auth"
	userUsecase "notes-app/internal/usecase/user"
)

// NewMux は依存関係を組み立て、ルートを登録した http.Handler を返す。
func NewMux(db *sql.DB) http.Handler {
	repo := userRepo.NewUser(db)
	userUC := userUsecase.NewUserUseCase(repo)
	userH := userHandler.New(userUC)
	loginUC := authUsecase.NewLoginUseCase(repo)
	loginH := authHandler.NewLoginHandler(loginUC)

	mux := http.NewServeMux()

	registerRootRoute(mux)
	registerHealthRoutes(mux)
	registerUserRoutes(mux, userH)
	registerAuthRoutes(mux, loginH)

	return mux
}
