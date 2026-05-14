package auth

import (
	"context"
	"errors"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"

	domain "notes-app/internal/domain/user"
	"notes-app/internal/models"

	useruc "notes-app/internal/usecase/user"
)

func seedLoginUser(t *testing.T, repo *useruc.StubRepo, email, plainPassword string) models.User {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("bcrypt: %v", err)
	}
	u := models.User{
		ID:           "11111111-1111-4111-8111-111111111111",
		NickName:     "login_tester",
		Email:        email,
		PasswordHash: string(hash),
		IsActive:     true,
	}
	created, err := repo.Create(context.Background(), u)
	if err != nil {
		t.Fatalf("stub Create: %v", err)
	}
	return created
}

func TestLoginUseCase_Execute_Success(t *testing.T) {
	t.Setenv("JWT_SECRET", "unit-test-jwt-secret")
	repo := &useruc.StubRepo{}
	seedLoginUser(t, repo, "ok@example.com", "Password1!")

	uc := NewLoginUseCase(repo)
	token, err := uc.Execute(context.Background(), LoginInput{
		Email:    "ok@example.com",
		Password: "Password1!",
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
}

func TestLoginUseCase_Execute_UserNotFound(t *testing.T) {
	t.Setenv("JWT_SECRET", "unit-test-jwt-secret")
	repo := &useruc.StubRepo{}
	uc := NewLoginUseCase(repo)

	_, err := uc.Execute(context.Background(), LoginInput{
		Email:    "missing@example.com",
		Password: "Password1!",
	})
	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf("want ErrUserNotFound, got %v", err)
	}
}

func TestLoginUseCase_Execute_InvalidLogin_WrongPassword(t *testing.T) {
	t.Setenv("JWT_SECRET", "unit-test-jwt-secret")
	repo := &useruc.StubRepo{}
	seedLoginUser(t, repo, "u@example.com", "RightPass1!")
	uc := NewLoginUseCase(repo)

	_, err := uc.Execute(context.Background(), LoginInput{
		Email:    "u@example.com",
		Password: "WrongPass1!",
	})
	if !errors.Is(err, domain.ErrInvalidLogin) {
		t.Fatalf("want ErrInvalidLogin, got %v", err)
	}
}

func TestLoginUseCase_Execute_InvalidLogin_EmptyPassword(t *testing.T) {
	t.Setenv("JWT_SECRET", "unit-test-jwt-secret")
	repo := &useruc.StubRepo{}
	seedLoginUser(t, repo, "empty-pw@example.com", "Password1!")
	uc := NewLoginUseCase(repo)

	_, err := uc.Execute(context.Background(), LoginInput{
		Email:    "empty-pw@example.com",
		Password: "",
	})
	if !errors.Is(err, domain.ErrInvalidLogin) {
		t.Fatalf("want ErrInvalidLogin, got %v", err)
	}
}

func TestLoginUseCase_Execute_UserNotFound_LongEmail(t *testing.T) {
	t.Setenv("JWT_SECRET", "unit-test-jwt-secret")
	repo := &useruc.StubRepo{}
	uc := NewLoginUseCase(repo)

	longEmail := strings.Repeat("a", 200) + "@example.com"
	_, err := uc.Execute(context.Background(), LoginInput{
		Email:    longEmail,
		Password: "Password1!",
	})
	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf("want ErrUserNotFound, got %v", err)
	}
}
