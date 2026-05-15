package note

import (
	"log"
	"net/http"
	"time"

	"notes-app/internal/handler/httputil"
	"notes-app/internal/middleware"
	"notes-app/internal/usecase/note"
)

type NoteHandler struct {
	usecase *note.NoteUseCase
}

type NoteResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func New(usecase *note.NoteUseCase) *NoteHandler {
	return &NoteHandler{
		usecase: usecase,
	}
}

func (h *NoteHandler) Index(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	notes, err := h.usecase.Index(r.Context(), userID)
	if err != nil {
		log.Printf("note handler index: userID=%s err=%v", userID, err)
		httputil.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	var res []NoteResponse
	for _, n := range notes {
		res = append(res, NoteResponse{
			ID:        n.ID,
			Title:     n.Title,
			Content:   n.Content,
			CreatedAt: n.CreatedAt,
			UpdatedAt: n.UpdatedAt,
		})
	}

	httputil.WriteJSON(w, http.StatusOK, res)
}
