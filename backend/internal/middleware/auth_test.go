package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"notes-app/internal/handler/httputil"
)

const testJWTSecret = "middleware-test-secret"

func requireUnauthorizedJSONBody(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()
	var eb httputil.ErrorBody
	if err := json.Unmarshal(w.Body.Bytes(), &eb); err != nil {
		t.Fatalf("json: %v, body=%q", err, w.Body.String())
	}
	if eb.Message != clientAuthErrorMessage {
		t.Fatalf("message: want %q, got %q", clientAuthErrorMessage, eb.Message)
	}
}

func signedToken(t *testing.T, claims jwt.MapClaims) string {
	t.Helper()
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := tok.SignedString([]byte(testJWTSecret))
	if err != nil {
		t.Fatalf("SignedString: %v", err)
	}
	return s
}

func TestRequireAuth_Success(t *testing.T) {
	mw := NewAuthMiddleware(testJWTSecret)
	token := signedToken(t, jwt.MapClaims{
		"user_id": "user-uuid-123",
		"exp":     time.Now().Add(time.Hour).Unix(),
	})

	var gotUserID string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := r.Context().Value(UserIDKey)
		s, ok := v.(string)
		if !ok {
			t.Fatalf("context user_id: want string, got %T", v)
		}
		gotUserID = s
		w.WriteHeader(http.StatusTeapot)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	mw.RequireAuth(next).ServeHTTP(w, req)

	if w.Code != http.StatusTeapot {
		t.Fatalf("status: want %d, got %d, body=%q", http.StatusTeapot, w.Code, w.Body.String())
	}
	if gotUserID != "user-uuid-123" {
		t.Fatalf("user id: want %q, got %q", "user-uuid-123", gotUserID)
	}
}

func TestRequireAuth_MissingHeader(t *testing.T) {
	mw := NewAuthMiddleware(testJWTSecret)
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next must not be called")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	mw.RequireAuth(next).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status: want %d, got %d", http.StatusUnauthorized, w.Code)
	}
	requireUnauthorizedJSONBody(t, w)
}

func TestRequireAuth_InvalidBearerPrefix(t *testing.T) {
	mw := NewAuthMiddleware(testJWTSecret)
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next must not be called")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Basic xyz")
	w := httptest.NewRecorder()
	mw.RequireAuth(next).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status: want %d, got %d", http.StatusUnauthorized, w.Code)
	}
	requireUnauthorizedJSONBody(t, w)
}

func TestRequireAuth_LowercaseBearerRejected(t *testing.T) {
	mw := NewAuthMiddleware(testJWTSecret)
	token := signedToken(t, jwt.MapClaims{
		"user_id": "u1",
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next must not be called")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "bearer "+token)
	w := httptest.NewRecorder()
	mw.RequireAuth(next).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status: want %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestRequireAuth_MalformedToken(t *testing.T) {
	mw := NewAuthMiddleware(testJWTSecret)
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next must not be called")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer not-a-jwt")
	w := httptest.NewRecorder()
	mw.RequireAuth(next).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status: want %d, got %d", http.StatusUnauthorized, w.Code)
	}
	requireUnauthorizedJSONBody(t, w)
}

func TestRequireAuth_WrongSecret(t *testing.T) {
	wrongTok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "u1",
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	s, err := wrongTok.SignedString([]byte("other-secret"))
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	mw := NewAuthMiddleware(testJWTSecret)
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next must not be called")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+s)
	w := httptest.NewRecorder()
	mw.RequireAuth(next).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status: want %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestRequireAuth_ExpiredToken(t *testing.T) {
	mw := NewAuthMiddleware(testJWTSecret)
	token := signedToken(t, jwt.MapClaims{
		"user_id": "u1",
		"exp":     time.Now().Add(-2 * time.Hour).Unix(),
	})
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next must not be called")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	mw.RequireAuth(next).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status: want %d, got %d, body=%q", http.StatusUnauthorized, w.Code, w.Body.String())
	}
}

func TestRequireAuth_MissingUserIDClaim(t *testing.T) {
	mw := NewAuthMiddleware(testJWTSecret)
	token := signedToken(t, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next must not be called")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	mw.RequireAuth(next).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status: want %d, got %d", http.StatusUnauthorized, w.Code)
	}
	requireUnauthorizedJSONBody(t, w)
}

func TestRequireAuth_EmptyUserIDString(t *testing.T) {
	mw := NewAuthMiddleware(testJWTSecret)
	token := signedToken(t, jwt.MapClaims{
		"user_id": "",
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next must not be called")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	mw.RequireAuth(next).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status: want %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestRequireAuth_UserIDWrongJSONType(t *testing.T) {
	mw := NewAuthMiddleware(testJWTSecret)
	token := signedToken(t, jwt.MapClaims{
		"user_id": float64(999),
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next must not be called")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	mw.RequireAuth(next).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status: want %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestRequireAuth_EmptyBearerToken(t *testing.T) {
	mw := NewAuthMiddleware(testJWTSecret)
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next must not be called")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer ")
	w := httptest.NewRecorder()
	mw.RequireAuth(next).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status: want %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestRequireAuth_VeryLongAuthorizationHeader(t *testing.T) {
	mw := NewAuthMiddleware(testJWTSecret)
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next must not be called")
	})

	long := "Bearer " + strings.Repeat("x", 1<<12)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", long)
	w := httptest.NewRecorder()
	mw.RequireAuth(next).ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status: want %d, got %d", http.StatusUnauthorized, w.Code)
	}
}
