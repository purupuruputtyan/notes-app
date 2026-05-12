package user

import (
	"encoding/json"
	"net/http"

	"notes-app/internal/usecase/user"
)

type UserHandler struct {
	usecase *user.UserUseCase
}

type CreateUserRequest struct {
	NickName  string `json:"nick_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	IconImage string `json:"icon_image"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func New(usecase *user.UserUseCase) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	input := user.CreateUserInput{
		NickName:  req.NickName,
		Email:     req.Email,
		Password:  req.Password,
		IconImage: req.IconImage,
	}

	created, err := h.usecase.Create(input)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Message: err.Error(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(created); err != nil {
		http.Error(w, "failed to encode json", http.StatusInternalServerError)
		return
	}
}
