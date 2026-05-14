package httputil

import (
	"encoding/json"
	"errors"
	"net/http"
)

// DecodeJSON は r.Body を maxBytes まで読み dst にデコードする。
// 失敗時は w に応答を書き込み、非 nil の err を返す（呼び出し側は return するだけでよい）。
func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any, maxBytes int64) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			http.Error(w, "request too large", http.StatusRequestEntityTooLarge)
			return err
		}
		http.Error(w, "invalid request", http.StatusBadRequest)
		return err
	}
	return nil
}
