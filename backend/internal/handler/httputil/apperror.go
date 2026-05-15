package httputil

import (
	"errors"
	"net/http"

	"notes-app/internal/apperror"
)

// ClientStatusFromAppError は apperror を HTTP ステータスとクライアント向けメッセージに落とす。
// 内部エラーはメッセージを伏せる。
func ClientStatusFromAppError(err error) (status int, message string) {
	if errors.Is(err, apperror.ErrUserNotFound) || errors.Is(err, apperror.ErrOwnerNotFound) {
		return http.StatusNotFound, "not found"
	}
	if ae, ok := errors.AsType[apperror.AppError](err); ok {
		return http.StatusBadRequest, ae.Message
	}
	return http.StatusInternalServerError, "internal server error"
}
