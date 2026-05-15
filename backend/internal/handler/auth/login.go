package auth

import (
	"net/http"

	"notes-app/internal/handler/httputil"
	usecase "notes-app/internal/usecase/auth"
)

type LoginHandler struct {
	usecase *usecase.LoginUseCase
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func NewLoginHandler(uc *usecase.LoginUseCase) *LoginHandler {
	return &LoginHandler{usecase: uc}
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req LoginRequest
	if err := httputil.DecodeJSON(w, r, &req, httputil.DefaultMaxJSONBody); err != nil {
		return
	}

	token, err := h.usecase.Execute(r.Context(), usecase.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		status, msg := httputil.ClientStatusFromAppError(err)
		httputil.WriteError(w, status, msg)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, LoginResponse{Token: token})
}
