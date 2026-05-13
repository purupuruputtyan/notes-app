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
	for _, u := range s.users {
		if u.ID == id {
			return u, nil
		}
	}
	return models.User{}, domain.ErrUserNotFound
}

func (s *stubRepo) Update(id string, params domain.UpdateUserParams) (models.User, error) {
	for i, u := range s.users {
		if u.ID == id {
			s.users[i].NickName = params.NickName
			s.users[i].Email = params.Email
			s.users[i].PasswordHash = params.PasswordHash
			s.users[i].IconImage = params.IconImage
			return s.users[i], nil
		}
	}
	return models.User{}, domain.ErrUserNotFound
}

func setupHandler() *UserHandler {
	repo := &stubRepo{}
	uc := user.NewUserUseCase(repo)
	return New(uc)
}

func newRequest(method, url, body string) *http.Request {
	req := httptest.NewRequest(method, url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func putUserRequest(userID, body string) *http.Request {
	return newRequest(http.MethodPut, "/users/"+userID, body)
}

func defaultUserRequest() UserRequest {
	return UserRequest{
		NickName:  "テストユーザー",
		Email:     "test@example.com",
		Password:  "password123!",
		IconImage: "https://example.com",
	}
}

func mustMarshalUserRequest(t *testing.T, req UserRequest) string {
	t.Helper()
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	return string(b)
}

func assertHTTPStatus(t *testing.T, w *httptest.ResponseRecorder, want int) {
	t.Helper()
	if w.Code != want {
		t.Fatalf("expected status %d, got %d", want, w.Code)
	}
}

func decodeUser(t *testing.T, w *httptest.ResponseRecorder) models.User {
	t.Helper()
	var u models.User
	if err := json.Unmarshal(w.Body.Bytes(), &u); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	return u
}

func decodeShowUserResponse(t *testing.T, w *httptest.ResponseRecorder) ShowUserResponse {
	t.Helper()
	var got ShowUserResponse
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode show response: %v", err)
	}
	return got
}

func createUser(t *testing.T, h *UserHandler) models.User {
	t.Helper()

	u, err := h.usecase.Create(user.CreateUserInput{
		NickName:  "テストユーザー",
		Email:     "test@example.com",
		Password:  "Password123!",
		IconImage: "https://example.com",
	})
	if err != nil {
		t.Fatalf("create error: %v", err)
	}
	return u
}

func TestUserHandler_Create(t *testing.T) {
	h := setupHandler()
	w := httptest.NewRecorder()
	body := mustMarshalUserRequest(t, UserRequest{
		NickName:  "テストユーザー",
		Email:     "test@example.com",
		Password:  "password123!",
		IconImage: "https://example.com",
	})

	h.Create(w, newRequest(http.MethodPost, "/users", body))

	assertHTTPStatus(t, w, http.StatusCreated)

	got := decodeUser(t, w)

	if got.ID == "" {
		t.Fatalf("expected ID to be set")
	}

	if got.PasswordHash == "password123!" {
		t.Fatalf("expected password to be hashed")
	}
}

func TestUserHandler_Create_Validation(t *testing.T) {
	const invalidJSONBody = `{
			"nick_name":"テストユーザー",
			"email":"test@example.com",
			"password":"password123!",
			"icon_image":"https://example.co
		`

	cases := []struct {
		name    string
		mutate  func(*UserRequest)
		rawBody string
	}{
		{
			name: "EmptyNickName",
			mutate: func(r *UserRequest) {
				r.NickName = ""
			},
		},
		{
			name: "NickNameTooLong",
			mutate: func(r *UserRequest) {
				r.NickName = strings.Repeat("a", 101)
			},
		},
		{
			name: "InvalidEmail",
			mutate: func(r *UserRequest) {
				r.Email = "testexample.com"
			},
		},
		{
			name: "PasswordTooShort",
			mutate: func(r *UserRequest) {
				r.Password = "pa123!"
			},
		},
		{
			name: "InvalidPassword",
			mutate: func(r *UserRequest) {
				r.Password = "password"
			},
		},
		{
			name:    "InvalidJSON",
			rawBody: invalidJSONBody,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			h := setupHandler()
			w := httptest.NewRecorder()

			var body string
			switch {
			case tt.rawBody != "":
				body = tt.rawBody
			default:
				req := defaultUserRequest()
				tt.mutate(&req)
				body = mustMarshalUserRequest(t, req)
			}

			h.Create(w, newRequest(http.MethodPost, "/users", body))
			assertHTTPStatus(t, w, http.StatusBadRequest)
		})
	}
}

func TestUserHandler_Show(t *testing.T) {
	h := setupHandler()
	created := createUser(t, h)

	req := httptest.NewRequest(
		http.MethodGet,
		"/users/"+created.ID,
		nil,
	)

	w := httptest.NewRecorder()

	h.Show(w, req, created.ID)

	assertHTTPStatus(t, w, http.StatusOK)

	got := decodeShowUserResponse(t, w)

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

	assertHTTPStatus(t, w, http.StatusNotFound)

	if !strings.Contains(w.Body.String(), "not found") {
		t.Fatalf(
			"expected response body to contain not found",
		)
	}
}

func TestUserHandler_Update(t *testing.T) {
	h := setupHandler()
	created := createUser(t, h)

	w := httptest.NewRecorder()
	body := mustMarshalUserRequest(t, UserRequest{
		NickName:  "アップデートユーザー",
		Email:     "update@example.com",
		Password:  "updated123!",
		IconImage: "https://update.com",
	})

	h.Update(w, putUserRequest(created.ID, body), created.ID)

	assertHTTPStatus(t, w, http.StatusOK)

	got := decodeUser(t, w)

	if got.NickName != "アップデートユーザー" {
		t.Fatalf("nick mismatch: %s", got.NickName)
	}

	if got.Email != "update@example.com" {
		t.Fatalf("email mismatch: %s", got.Email)
	}

	if got.IconImage.String != "https://update.com" {
		t.Fatalf("icon mismatch: %v", got.IconImage)
	}

	if got.PasswordHash == "updated123!" {
		t.Fatal("password not hashed")
	}
}

func TestUserHandler_Update_Validation(t *testing.T) {
	cases := []struct {
		name    string
		mutate  func(*UserRequest)
		rawBody string
	}{
		{
			name: "EmptyNickName",
			mutate: func(r *UserRequest) {
				r.NickName = ""
			},
		},
		{
			name: "NickNameTooLong",
			mutate: func(r *UserRequest) {
				r.NickName = strings.Repeat("a", 101)
			},
		},
		{
			name: "InvalidEmail",
			mutate: func(r *UserRequest) {
				r.Email = "testexample.com"
			},
		},
		{
			name: "PasswordTooShort",
			mutate: func(r *UserRequest) {
				r.Password = "pa123!"
			},
		},
		{
			name: "InvalidPassword",
			mutate: func(r *UserRequest) {
				r.Password = "password"
			},
		},
		{
			name:    "InvalidJSON",
			rawBody: `{"nick_name":"テストユーザー"`,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			h := setupHandler()
			created := createUser(t, h)
			w := httptest.NewRecorder()

			var body string
			switch {
			case tt.rawBody != "":
				body = tt.rawBody
			default:
				req := defaultUserRequest()
				tt.mutate(&req)
				body = mustMarshalUserRequest(t, req)
			}

			h.Update(w, putUserRequest(created.ID, body), created.ID)
			assertHTTPStatus(t, w, http.StatusBadRequest)
		})
	}
}
