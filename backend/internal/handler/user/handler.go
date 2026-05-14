package user

import (
	"encoding/json"
	"errors"
	"net/http"

	domain "notes-app/internal/domain/user"
	"notes-app/internal/usecase/user"
)

// リクエスト JSON の上限（DoS 対策）
const maxJSONBodyBytes = 1 << 20 // 1 MiB

type UserHandler struct {
	usecase *user.UserUseCase
}

type UserRequest struct {
	NickName  string `json:"nick_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	IconImage string `json:"icon_image"`
}

type ShowUserResponse struct {
	ID        string `json:"id"`
	NickName  string `json:"nick_name"`
	Email     string `json:"email"`
	IconImage string `json:"icon_image"`
	IsActive  bool   `json:"is_active"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func New(usecase *user.UserUseCase) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}

// writeJSON は v を JSON 化してから書き込む（WriteHeader 後の Encode 失敗を避ける）。
func writeJSON(w http.ResponseWriter, status int, v any) {
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

func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Message: message})
}

func decodeUserJSON(w http.ResponseWriter, r *http.Request, out *UserRequest) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBodyBytes)
	if err := json.NewDecoder(r.Body).Decode(out); err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			http.Error(w, "request too large", http.StatusRequestEntityTooLarge)
			return false
		}
		http.Error(w, "invalid request", http.StatusBadRequest)
		return false
	}
	return true
}

// mapUsecaseError はクライアント向けステータスとメッセージに落とす（内部エラーは伏せる）。
func mapUsecaseError(err error) (status int, message string) {
	if errors.Is(err, domain.ErrUserNotFound) {
		return http.StatusNotFound, "not found"
	}
	var ae domain.AppError
	if errors.As(err, &ae) {
		return http.StatusBadRequest, ae.Message
	}
	return http.StatusInternalServerError, "internal server error"
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req UserRequest
	if !decodeUserJSON(w, r, &req) {
		return
	}

	input := user.CreateUserInput{
		NickName:  req.NickName,
		Email:     req.Email,
		Password:  req.Password,
		IconImage: req.IconImage,
	}

	created, err := h.usecase.Create(r.Context(), input)
	if err != nil {
		status, msg := mapUsecaseError(err)
		writeJSONError(w, status, msg)
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *UserHandler) Show(w http.ResponseWriter, r *http.Request, id string) {
	found, err := h.usecase.Show(r.Context(), id)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "not found")
		return
	}

	res := ShowUserResponse{
		ID:        found.ID,
		NickName:  found.NickName,
		Email:     found.Email,
		IconImage: found.IconImage.String,
		IsActive:  found.IsActive,
	}

	writeJSON(w, http.StatusOK, res)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request, id string) {
	defer r.Body.Close()

	var req UserRequest
	if !decodeUserJSON(w, r, &req) {
		return
	}

	input := user.UpdateUserInput{
		ID:        id,
		NickName:  req.NickName,
		Email:     req.Email,
		Password:  req.Password,
		IconImage: req.IconImage,
	}

	updated, err := h.usecase.Update(r.Context(), input)
	if err != nil {
		status, msg := mapUsecaseError(err)
		writeJSONError(w, status, msg)
		return
	}

	writeJSON(w, http.StatusOK, updated)
}
