// Package httputil は internal/handler 配下で共有する HTTP / JSON の小さなヘルパを提供する。
package httputil

import (
	"encoding/json"
	"net/http"
)

// DefaultMaxJSONBody は JSON リクエストボディの既定上限バイト数（DoS 対策）。
const DefaultMaxJSONBody int64 = 1 << 20 // 1 MiB

// ErrorBody はクライアント向けの単純な JSON エラー本文。
type ErrorBody struct {
	Message string `json:"message"`
}

// WriteJSON は v を JSON 化してから書き込む（WriteHeader 後の Encode 失敗を避ける）。
func WriteJSON(w http.ResponseWriter, status int, v any) {
	b, err := json.Marshal(v)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("failed to encode json"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(b)
}

// WriteError は JSON で { "message": ... } を返す。
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, ErrorBody{Message: message})
}
