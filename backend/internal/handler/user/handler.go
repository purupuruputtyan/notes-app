package user

import (
	"encoding/json"
	"net/http"

	"notes-app/internal/usecase/user"
)

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

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req UserRequest

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

	created, err := h.usecase.Create(r.Context(), input)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(created); err != nil {
		http.Error(w, "failed to encode json", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) Show(w http.ResponseWriter, r *http.Request, id string) {
	found, err := h.usecase.Show(r.Context(), id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	res := ShowUserResponse{
		ID:        found.ID,
		NickName:  found.NickName,
		Email:     found.Email,
		IconImage: found.IconImage.String,
		IsActive:  found.IsActive,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "failed to encode json", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request, id string) {
	var req UserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(updated); err != nil {
		http.Error(w, "failed to encode json", http.StatusInternalServerError)
		return
	}
}
