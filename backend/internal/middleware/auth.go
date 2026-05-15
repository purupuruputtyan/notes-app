package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"notes-app/internal/handler/httputil"
)

type contextKey string

const UserIDKey contextKey = "user_id"

// clientAuthErrorMessage は認証失敗時にクライアントへ返す一律メッセージ（詳細はログに出さない）。
const clientAuthErrorMessage = "unauthorized"

type AuthMiddleware struct {
	jwtSecret string
}

func NewAuthMiddleware(secret string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: secret,
	}
}

func (m *AuthMiddleware) RequireAuth(
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			httputil.WriteError(w, http.StatusUnauthorized, clientAuthErrorMessage)
			return
		}

		const bearerPrefix = "Bearer "

		if !strings.HasPrefix(authHeader, bearerPrefix) {
			httputil.WriteError(w, http.StatusUnauthorized, clientAuthErrorMessage)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}

			return []byte(m.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			httputil.WriteError(w, http.StatusUnauthorized, clientAuthErrorMessage)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			httputil.WriteError(w, http.StatusUnauthorized, clientAuthErrorMessage)
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			httputil.WriteError(w, http.StatusUnauthorized, clientAuthErrorMessage)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
