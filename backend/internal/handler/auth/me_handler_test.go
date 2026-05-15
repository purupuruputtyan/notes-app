package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aarondl/null/v8"

	"notes-app/internal/handler/httputil"
	"notes-app/internal/middleware"
	"notes-app/internal/models"
	authUc "notes-app/internal/usecase/auth"
	userUc "notes-app/internal/usecase/user"
)

// boomShowRepo は Show のみ失敗させ、内部エラー経路を検証する。
type boomShowRepo struct {
	*userUc.StubRepo
}

func (b *boomShowRepo) Show(ctx context.Context, id string) (models.User, error) {
	_ = ctx
	_ = id
	return models.User{}, fmt.Errorf("simulated repository failure")
}

func meHandlerTestSetup(t *testing.T) (*MeHandler, *userUc.StubRepo) {
	t.Helper()
	repo := &userUc.StubRepo{}
	h := NewMeHandler(authUc.NewMeUseCase(repo))
	return h, repo
}

func newMeRequestWithUserID(t *testing.T, userID string) *http.Request {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	if userID != "" {
		ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
		req = req.WithContext(ctx)
	}
	return req
}

func TestMeHandler_Me_Success(t *testing.T) {
	h, repo := meHandlerTestSetup(t)
	ctx := context.Background()
	_, err := repo.Create(ctx, models.User{
		ID:           "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa",
		NickName:     "nick",
		Email:        "nick@example.com",
		PasswordHash: "hash",
		IconImage:    null.StringFrom("https://example.com/icon.png"),
		IsActive:     true,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	w := httptest.NewRecorder()
	h.Me(w, newMeRequestWithUserID(t, "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa"))

	if w.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d, body=%q", http.StatusOK, w.Code, w.Body.String())
	}
	var got MeResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("json: %v", err)
	}
	if got.ID != "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa" {
		t.Fatalf("id: want uuid, got %q", got.ID)
	}
	if got.NickName != "nick" {
		t.Fatalf("nick_name: want nick, got %q", got.NickName)
	}
	if got.Email != "nick@example.com" {
		t.Fatalf("email: %q", got.Email)
	}
	if got.IconImage != "https://example.com/icon.png" {
		t.Fatalf("icon_image: %q", got.IconImage)
	}
	if !got.IsActive {
		t.Fatal("is_active: want true")
	}
}

func TestMeHandler_Me_Success_EmptyIconImage(t *testing.T) {
	h, repo := meHandlerTestSetup(t)
	ctx := context.Background()
	_, err := repo.Create(ctx, models.User{
		ID:           "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb",
		NickName:     "no_icon",
		Email:        "no-icon@example.com",
		PasswordHash: "hash",
		IconImage:    null.String{},
		IsActive:     false,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	w := httptest.NewRecorder()
	h.Me(w, newMeRequestWithUserID(t, "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb"))

	if w.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d", http.StatusOK, w.Code)
	}
	var got MeResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("json: %v", err)
	}
	if got.IconImage != "" {
		t.Fatalf("icon_image: want empty, got %q", got.IconImage)
	}
	if got.IsActive {
		t.Fatal("is_active: want false")
	}
}

func TestMeHandler_Me_UnauthorizedMissingContext(t *testing.T) {
	h, _ := meHandlerTestSetup(t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	h.Me(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status: want %d, got %d, body=%q", http.StatusUnauthorized, w.Code, w.Body.String())
	}
	var eb httputil.ErrorBody
	if err := json.Unmarshal(w.Body.Bytes(), &eb); err != nil {
		t.Fatalf("json: %v", err)
	}
	if eb.Message == "" {
		t.Fatal("expected error message")
	}
}

func TestMeHandler_Me_NotFound(t *testing.T) {
	h, _ := meHandlerTestSetup(t)
	w := httptest.NewRecorder()
	h.Me(w, newMeRequestWithUserID(t, "cccccccc-cccc-4ccc-8ccc-cccccccccccc"))

	if w.Code != http.StatusNotFound {
		t.Fatalf("status: want %d, got %d, body=%q", http.StatusNotFound, w.Code, w.Body.String())
	}
}

func TestMeHandler_Me_InternalError(t *testing.T) {
	repo := &boomShowRepo{StubRepo: &userUc.StubRepo{}}
	h := NewMeHandler(authUc.NewMeUseCase(repo))
	w := httptest.NewRecorder()
	h.Me(w, newMeRequestWithUserID(t, "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa"))

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

func TestMeHandler_Me_WrongContextValueType(t *testing.T) {
	h, _ := meHandlerTestSetup(t)
	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserIDKey, 999))
	w := httptest.NewRecorder()
	h.Me(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status: want %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestMeHandler_Me_VeryLongUserIDNotFound(t *testing.T) {
	h, _ := meHandlerTestSetup(t)
	long := strings.Repeat("x", 2048)
	w := httptest.NewRecorder()
	h.Me(w, newMeRequestWithUserID(t, long))

	if w.Code != http.StatusNotFound {
		t.Fatalf("status: want %d, got %d", http.StatusNotFound, w.Code)
	}
}
