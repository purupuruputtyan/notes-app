package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	domain "notes-app/internal/domain/user"
	"notes-app/internal/models"
	"notes-app/internal/usecase/user"
)

type stubRepo struct {
	users []models.User
}

func (s *stubRepo) Create(u models.User) (models.User, error) {
	u.ID = "test-id"
	s.users = append(s.users, u)
	return u, nil
}

func (s *stubRepo) Show(id string) (models.User, error) {
	for _, user := range s.users {
		if user.ID == id {
			return user, nil
		}
	}

	return models.User{}, domain.ErrUserNotFound
}

func setupHandler() *UserHandler {
	repo := &stubRepo{}
	uc := user.NewUserUseCase(repo)
	return New(uc)
}

func TestUserHandler_Create(t *testing.T) {
	h := setupHandler()

	reqBody := strings.NewReader(`{
		"nick_name":"テストユーザー",
		"email":"test@example.com",
		"password":"password123!",
		"icon_image":"https://example.com"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/users", reqBody)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	var got models.User
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if got.ID == "" {
		t.Fatalf("expected ID to be set")
	}

	if got.PasswordHash == "password123!" {
		t.Fatalf("expected password to be hashed")
	}
}

func TestUserHandler_Create_EmptyNickName(t *testing.T) {
	h := setupHandler()

	reqBody := strings.NewReader(`{
		"nick_name":"",
		"email":"test@example.com",
		"password":"password123!",
		"icon_image":"https://example.com"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/users", reqBody)
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUserHandler_Create_NickNameTooLong(t *testing.T) {
	h := setupHandler()

	longNickName := "a"
	for len(longNickName) <= 100 {
		longNickName += "a"
	}

	reqBody := strings.NewReader(`{
		"nick_name":"` + longNickName + `",
		"email":"test@example.com",
		"password":"password123!",
		"icon_image":"https://example.com"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/users", reqBody)
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUserHandler_Create_InvalidEmail(t *testing.T) {
	h := setupHandler()

	reqBody := strings.NewReader(`{
		"nick_name":"テストユーザー",
		"email":"testexample.com",
		"password":"password123!",
		"icon_image":"https://example.com"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/users", reqBody)
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUserHandler_Create_PasswordTooShort(t *testing.T) {
	h := setupHandler()

	reqBody := strings.NewReader(`{
		"nick_name":"テストユーザー",
		"email":"test@example.com",
		"password":"pa123!",
		"icon_image":"https://example.com"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/users", reqBody)
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUserHandler_Create_InvalidPassword(t *testing.T) {
	h := setupHandler()

	reqBody := strings.NewReader(`{
		"nick_name":"テストユーザー",
		"email":"test@example.com",
		"password":"password",
		"icon_image":"https://example.com"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/users", reqBody)
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUserHandler_Create_InvalidJSON(t *testing.T) {
	h := setupHandler()

	reqBody := strings.NewReader(`{
		"nick_name":"テストユーザー",
		"email":"test@example.com",
		"password":"password123!",
		"icon_image":"https://example.co
	`)

	req := httptest.NewRequest(http.MethodPost, "/users", reqBody)
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUserHandler_Show(t *testing.T) {
	h := setupHandler()

	created, _ := h.usecase.Create(user.CreateUserInput{
		NickName:  "テストユーザー",
		Email:     "test@example.com",
		Password:  "Password123!",
		IconImage: "https://example.com",
	})

	req := httptest.NewRequest(
		http.MethodGet,
		"/users/"+created.ID,
		nil,
	)

	w := httptest.NewRecorder()

	h.Show(w, req, created.ID)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var got ShowUserResponse

	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.ID != created.ID {
		t.Fatalf("expected id %s, got %s", created.ID, got.ID)
	}

	if got.NickName != created.NickName {
		t.Fatalf(
			"expected nickname %s, got %s",
			created.NickName,
			got.NickName,
		)
	}

	if strings.Contains(w.Body.String(), "password_hash") {
		t.Fatalf("password_hash should not be returned")
	}
}

func TestUserHandler_Show_NotFound(t *testing.T) {
	h := setupHandler()

	req := httptest.NewRequest(
		http.MethodGet,
		"/users/not-found-id",
		nil,
	)

	w := httptest.NewRecorder()

	h.Show(w, req, "not-found-id")

	if w.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status 404, got %d",
			w.Code,
		)
	}

	if !strings.Contains(w.Body.String(), "not found") {
		t.Fatalf(
			"expected response body to contain not found",
		)
	}
}
