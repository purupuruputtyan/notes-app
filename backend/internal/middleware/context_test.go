package middleware

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetUserID_Present(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	ctx := context.WithValue(req.Context(), UserIDKey, "uid-1")
	req = req.WithContext(ctx)

	got, ok := GetUserID(req)
	if !ok {
		t.Fatal("expected ok true")
	}
	if got != "uid-1" {
		t.Fatalf("user id: want %q, got %q", "uid-1", got)
	}
}

func TestGetUserID_Missing(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)

	_, ok := GetUserID(req)
	if ok {
		t.Fatal("expected ok false when context has no user id")
	}
}

func TestGetUserID_WrongTypeInContext(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req = req.WithContext(context.WithValue(req.Context(), UserIDKey, 42))

	_, ok := GetUserID(req)
	if ok {
		t.Fatal("expected ok false when value is not a string")
	}
}

func TestGetUserID_EmptyStringStillOkTrue(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req = req.WithContext(context.WithValue(req.Context(), UserIDKey, ""))

	got, ok := GetUserID(req)
	if !ok {
		t.Fatal("empty string is still type string; ok should be true")
	}
	if got != "" {
		t.Fatalf("want empty string, got %q", got)
	}
}

func TestGetUserID_VeryLongUserID(t *testing.T) {
	long := strings.Repeat("a", 1<<12)

	req := httptest.NewRequest("GET", "/", nil)
	req = req.WithContext(context.WithValue(req.Context(), UserIDKey, long))

	got, ok := GetUserID(req)
	if !ok {
		t.Fatal("expected ok true")
	}
	if len(got) != 1<<12 {
		t.Fatalf("length: want %d, got %d", 1<<12, len(got))
	}
}
