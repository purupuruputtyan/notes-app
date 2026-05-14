package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"notes-app/internal/handler/httputil"
	"notes-app/internal/models"
	authuc "notes-app/internal/usecase/auth"
	useruc "notes-app/internal/usecase/user"
)

func loginSetupHandler(t *testing.T) *LoginHandler {
	t.Helper()
	t.Setenv("JWT_SECRET", "handler-test-jwt-secret")
	repo := &useruc.StubRepo{}
	hash, err := bcrypt.GenerateFromPassword([]byte("Password1!"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("bcrypt: %v", err)
	}
	_, err = repo.Create(context.Background(), models.User{
		ID:           "22222222-2222-4222-8222-222222222222",
		NickName:     "handler_login",
		Email:        "handler-login@example.com",
		PasswordHash: string(hash),
		IsActive:     true,
	})
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}
	return NewLoginHandler(authuc.NewLoginUseCase(repo))
}

func postLoginJSON(t *testing.T, body string) *http.Request {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestLoginHandler_Login_Success(t *testing.T) {
	h := loginSetupHandler(t)
	w := httptest.NewRecorder()
	body := `{"email":"handler-login@example.com","password":"Password1!"}`
	h.Login(w, postLoginJSON(t, body))

	if w.Code != http.StatusOK {
		t.Fatalf("status: want %d, got %d, body=%q", http.StatusOK, w.Code, w.Body.String())
	}
	var got LoginResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("json: %v", err)
	}
	if got.Token == "" {
		t.Fatal("expected token in response")
	}
}

func TestLoginHandler_Login_InvalidPassword(t *testing.T) {
	h := loginSetupHandler(t)
	w := httptest.NewRecorder()
	body := `{"email":"handler-login@example.com","password":"WrongPassword1!"}`
	h.Login(w, postLoginJSON(t, body))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: want %d, got %d, body=%q", http.StatusBadRequest, w.Code, w.Body.String())
	}
	var eb httputil.ErrorBody
	if err := json.Unmarshal(w.Body.Bytes(), &eb); err != nil {
		t.Fatalf("json: %v", err)
	}
	if eb.Message == "" {
		t.Fatal("expected error message")
	}
}

func TestLoginHandler_Login_UserNotFound(t *testing.T) {
	h := loginSetupHandler(t)
	w := httptest.NewRecorder()
	body := `{"email":"nobody@example.com","password":"Password1!"}`
	h.Login(w, postLoginJSON(t, body))

	if w.Code != http.StatusNotFound {
		t.Fatalf("status: want %d, got %d, body=%q", http.StatusNotFound, w.Code, w.Body.String())
	}
}

func TestLoginHandler_Login_InvalidJSON(t *testing.T) {
	h := loginSetupHandler(t)
	w := httptest.NewRecorder()
	h.Login(w, postLoginJSON(t, `{"email":`))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: want %d, got %d", http.StatusBadRequest, w.Code)
	}
	if !strings.Contains(w.Body.String(), "invalid request") {
		t.Fatalf("body=%q", w.Body.String())
	}
}

func TestLoginHandler_Login_EmptyBody(t *testing.T) {
	h := loginSetupHandler(t)
	w := httptest.NewRecorder()
	h.Login(w, postLoginJSON(t, ""))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: want %d, got %d", http.StatusBadRequest, w.Code)
	}
}
