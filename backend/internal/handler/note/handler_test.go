package note

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"notes-app/internal/handler/httputil"
	"notes-app/internal/middleware"
	"notes-app/internal/models"
	noteUc "notes-app/internal/usecase/note"
)

func setupHandler(repo *noteUc.StubRepo) *NoteHandler {
	uc := noteUc.NewNoteUseCase(repo)
	return New(uc)
}

// indexRequest は userID をコンテキストにセットした GET /notes リクエストを返す。
// 実際のルーティングでは authMiddleware が同じ方法で userID を注入するため、
// この関数はその動作を再現している。
func indexRequest(userID string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/notes", nil)
	if userID != "" {
		ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
		req = req.WithContext(ctx)
	}
	return req
}

func decodeNotesResponse(t *testing.T, w *httptest.ResponseRecorder) []NoteResponse {
	t.Helper()
	var got []NoteResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode notes: %v, body=%q", err, w.Body.String())
	}
	return got
}

func TestNoteHandler_Index_Success(t *testing.T) {
	userID := "11111111-1111-4111-8111-111111111111"
	ts := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	repo := noteUc.NewTestStubRepo(models.NoteSlice{
		{ID: "n1", UserID: userID, Title: "first", Content: "body1", CreatedAt: ts, UpdatedAt: ts},
		{ID: "n2", UserID: userID, Title: "second", Content: "body2", CreatedAt: ts, UpdatedAt: ts},
		{ID: "n3", UserID: "other-user", Title: "skip", Content: "x", CreatedAt: ts, UpdatedAt: ts},
	}, nil)
	h := setupHandler(repo)
	w := httptest.NewRecorder()

	h.Index(w, indexRequest(userID))

	if w.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d, body=%q", http.StatusOK, w.Code, w.Body.String())
	}
	got := decodeNotesResponse(t, w)
	if len(got) != 2 {
		t.Fatalf("len: want 2, got %d", len(got))
	}
	if got[0].ID != "n1" || got[0].Title != "first" || got[0].Content != "body1" {
		t.Fatalf("first note: %+v", got[0])
	}
	if got[1].Title != "second" {
		t.Fatalf("second note: %+v", got[1])
	}
	if !got[0].CreatedAt.Equal(ts) {
		t.Fatalf("CreatedAt: want %v, got %v", ts, got[0].CreatedAt)
	}
	if strings.Contains(w.Body.String(), "user_id") {
		t.Fatalf("user_id must not appear in response body: %s", w.Body.String())
	}
}

func TestNoteHandler_Index_Empty(t *testing.T) {
	userID := "22222222-2222-4222-8222-222222222222"
	repo := noteUc.NewTestStubRepo(models.NoteSlice{
		{ID: "n1", UserID: "other", Title: "t", Content: "c"},
	}, nil)
	h := setupHandler(repo)
	w := httptest.NewRecorder()

	h.Index(w, indexRequest(userID))

	if w.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d", http.StatusOK, w.Code)
	}
	got := decodeNotesResponse(t, w)
	if len(got) != 0 {
		t.Fatalf("want empty slice, got len=%d", len(got))
	}
}

// TestNoteHandler_Index_MissingUserIDContext は、認証ミドルウェアが
// コンテキストに userID を注入していない状態（認可バイパス等）を想定したテスト。
// 実運用では authMiddleware が先に 401 を返すため、このケースはハンドラまで到達しないが、
// 多重防衛として 401 を返すことを保証する。
func TestNoteHandler_Index_MissingUserIDContext(t *testing.T) {
	repo := noteUc.NewTestStubRepo(nil, nil)
	h := setupHandler(repo)
	w := httptest.NewRecorder()

	// コンテキストに userID をセットしないリクエスト
	h.Index(w, indexRequest(""))

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status: want %d, got %d, body=%q", http.StatusUnauthorized, w.Code, w.Body.String())
	}
	var eb httputil.ErrorBody
	if err := json.Unmarshal(w.Body.Bytes(), &eb); err != nil {
		t.Fatalf("json: %v", err)
	}
	if eb.Message != "unauthorized" {
		t.Fatalf("message: want %q, got %q", "unauthorized", eb.Message)
	}
}

func TestNoteHandler_Index_InternalError(t *testing.T) {
	userID := "11111111-1111-4111-8111-111111111111"
	dbErr := errors.New("db unavailable")
	repo := noteUc.NewTestStubRepo(nil, dbErr)
	h := setupHandler(repo)
	w := httptest.NewRecorder()

	h.Index(w, indexRequest(userID))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status: want %d, got %d, body=%q", http.StatusInternalServerError, w.Code, w.Body.String())
	}
	var eb httputil.ErrorBody
	if err := json.Unmarshal(w.Body.Bytes(), &eb); err != nil {
		t.Fatalf("json: %v", err)
	}
	if eb.Message != "internal server error" {
		t.Fatalf("message: want %q, got %q", "internal server error", eb.Message)
	}
}
