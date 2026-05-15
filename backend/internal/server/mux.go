package server

import (
	"database/sql"
	"net/http"

	authHandler "notes-app/internal/handler/auth"
	userHandler "notes-app/internal/handler/user"
	middleware "notes-app/internal/middleware"
	userRepo "notes-app/internal/repository/user"
	authUsecase "notes-app/internal/usecase/auth"
	userUsecase "notes-app/internal/usecase/user"
)

// NewMux は依存関係を組み立て、ルートを登録した http.Handler を返す。
// jwtSecret は署名検証用の共有秘密（空であってはならない）。起動時に config.Load 等で検証した値を渡すこと。
func NewMux(db *sql.DB, jwtSecret string) http.Handler {
	repo := userRepo.NewUser(db)
	userUC := userUsecase.NewUserUseCase(repo)
	userH := userHandler.New(userUC)
	loginUC := authUsecase.NewLoginUseCase(repo)
	loginH := authHandler.NewLoginHandler(loginUC)
	meUC := authUsecase.NewMeUseCase(repo)
	meH := authHandler.NewMeHandler(meUC)

	authMiddleware := middleware.NewAuthMiddleware(jwtSecret)

	mux := http.NewServeMux()

	registerRootRoute(mux)
	registerHealthRoutes(mux)
	registerUserRoutes(mux, userH)
	registerAuthRoutes(mux, loginH)
	registerMeRoutes(mux, meH, authMiddleware)

	return mux
}
