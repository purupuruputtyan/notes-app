package httputil

import (
	"errors"
	"net/http"

	domain "notes-app/internal/domain/user"
)

// ClientStatusFromUserDomain は domain/user のエラーを HTTP ステータスとクライアント向けメッセージに落とす。
// 内部エラーはメッセージを伏せる。
func ClientStatusFromUserDomain(err error) (status int, message string) {
	if errors.Is(err, domain.ErrUserNotFound) {
		return http.StatusNotFound, "not found"
	}
	var ae domain.AppError
	if errors.As(err, &ae) {
		return http.StatusBadRequest, ae.Message
	}
	return http.StatusInternalServerError, "internal server error"
}
