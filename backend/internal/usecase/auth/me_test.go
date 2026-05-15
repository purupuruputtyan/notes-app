package auth

import (
	"context"
	"errors"
	"strings"
	"testing"

	"notes-app/internal/apperror"
	"notes-app/internal/models"
	userUc "notes-app/internal/usecase/user"
)

func TestMeUseCase_Execute_Success(t *testing.T) {
	repo := &userUc.StubRepo{}
	uc := NewMeUseCase(repo)
	ctx := context.Background()

	u, err := repo.Create(ctx, models.User{
		ID:           "11111111-1111-4111-8111-111111111111",
		NickName:     "me_uc",
		Email:        "me-uc@example.com",
		PasswordHash: "hash",
		IsActive:     true,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	got, err := uc.Execute(ctx, u.ID)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if got.ID != u.ID {
		t.Fatalf("ID: want %q, got %q", u.ID, got.ID)
	}
	if got.NickName != u.NickName {
		t.Fatalf("NickName: want %q, got %q", u.NickName, got.NickName)
	}
	if got.Email != u.Email {
		t.Fatalf("Email: want %q, got %q", u.Email, got.Email)
	}
}

func TestMeUseCase_Execute_NotFound(t *testing.T) {
	repo := &userUc.StubRepo{}
	uc := NewMeUseCase(repo)

	_, err := uc.Execute(
		context.Background(),
		"00000000-0000-4000-8000-000000000000",
	)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, apperror.ErrUserNotFound) {
		t.Fatalf("want ErrUserNotFound, got %v", err)
	}
}

func TestMeUseCase_Execute_EmptyUserID(t *testing.T) {
	repo := &userUc.StubRepo{}
	uc := NewMeUseCase(repo)

	_, err := uc.Execute(context.Background(), "")
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, apperror.ErrUserNotFound) {
		t.Fatalf("want ErrUserNotFound, got %v", err)
	}
}

func TestMeUseCase_Execute_VeryLongUserID(t *testing.T) {
	repo := &userUc.StubRepo{}
	uc := NewMeUseCase(repo)
	longID := strings.Repeat("z", 2048)

	_, err := uc.Execute(context.Background(), longID)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, apperror.ErrUserNotFound) {
		t.Fatalf("want ErrUserNotFound, got %v", err)
	}
}
