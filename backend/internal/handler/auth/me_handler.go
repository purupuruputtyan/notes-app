package auth

import (
	"errors"
	"net/http"

	"notes-app/internal/apperror"
	"notes-app/internal/handler/httputil"
	"notes-app/internal/middleware"
	usecase "notes-app/internal/usecase/auth"
)

type MeHandler struct {
	usecase *usecase.MeUseCase
}

func NewMeHandler(
	uc *usecase.MeUseCase,
) *MeHandler {
	return &MeHandler{
		usecase: uc,
	}
}

type MeResponse struct {
	ID        string `json:"id"`
	NickName  string `json:"nick_name"`
	Email     string `json:"email"`
	IconImage string `json:"icon_image"`
	IsActive  bool   `json:"is_active"`
}

func (h *MeHandler) Me(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		httputil.WriteError(
			w,
			http.StatusUnauthorized,
			"unauthorized: missing user id",
		)
		return
	}

	user, err := h.usecase.Execute(
		r.Context(),
		userID,
	)
	if err != nil {
		if errors.Is(err, apperror.ErrUserNotFound) {
			httputil.WriteError(
				w,
				http.StatusNotFound,
				"error: user not found",
			)
			return
		}
		httputil.WriteError(
			w,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	res := MeResponse{
		ID:        user.ID,
		NickName:  user.NickName,
		Email:     user.Email,
		IconImage: user.IconImage.String,
		IsActive:  user.IsActive,
	}

	httputil.WriteJSON(w, http.StatusOK, res)
}
