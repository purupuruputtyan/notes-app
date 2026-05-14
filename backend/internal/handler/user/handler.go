package user

import (
	"net/http"

	"notes-app/internal/handler/httputil"
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

func New(usecase *user.UserUseCase) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req UserRequest
	if err := httputil.DecodeJSON(w, r, &req, httputil.DefaultMaxJSONBody); err != nil {
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
		status, msg := httputil.ClientStatusFromUserDomain(err)
		httputil.WriteError(w, status, msg)
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, created)
}

func (h *UserHandler) Show(w http.ResponseWriter, r *http.Request, id string) {
	found, err := h.usecase.Show(r.Context(), id)
	if err != nil {
		httputil.WriteError(w, http.StatusNotFound, "not found")
		return
	}

	res := ShowUserResponse{
		ID:        found.ID,
		NickName:  found.NickName,
		Email:     found.Email,
		IconImage: found.IconImage.String,
		IsActive:  found.IsActive,
	}

	httputil.WriteJSON(w, http.StatusOK, res)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request, id string) {
	defer r.Body.Close()

	var req UserRequest
	if err := httputil.DecodeJSON(w, r, &req, httputil.DefaultMaxJSONBody); err != nil {
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
		status, msg := httputil.ClientStatusFromUserDomain(err)
		httputil.WriteError(w, status, msg)
		return
	}

	httputil.WriteJSON(w, http.StatusOK, updated)
}
