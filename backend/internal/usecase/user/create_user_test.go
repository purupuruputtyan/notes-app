package user

import (
	"strings"
	"testing"

	domain "notes-app/internal/domain/user"
)

func TestUserUseCase_Create(t *testing.T) {
	repo := &stubRepo{}
	uc := NewUserUseCase(repo)

	input := CreateUserInput{
		NickName:  "テストユーザー",
		Email:     "test@example.com",
		Password:  "Password-hash123",
		IconImage: "https://example.com",
	}

	created, err := uc.Create(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if created.ID == "" {
		t.Fatalf("expected ID to be set")
	}

	if created.NickName != "テストユーザー" {
		t.Fatalf("expected NickName テストユーザー, got %s", created.NickName)
	}

	if created.Email != "test@example.com" {
		t.Fatalf("expected Email test@example.com, got %s", created.Email)
	}

	if created.PasswordHash == "" {
		t.Fatalf("expected Password to be set")
	}

	if created.PasswordHash == "Password-hash123" {
		t.Fatalf("Password must be hashed")
	}

	if len(repo.users) != 1 {
		t.Fatalf("expected 1 user in repo, got %d", len(repo.users))
	}
}

func TestUserUseCase_Create_EmptyNickName(t *testing.T) {
	repo := &stubRepo{}
	uc := NewUserUseCase(repo)

	input := CreateUserInput{
		NickName:  "",
		Email:     "test@example.com",
		Password:  "Password-hash123",
		IconImage: "https://example.com",
	}

	_, err := uc.Create(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != domain.ErrNickNameRequired {
		t.Fatalf("expected ErrNickNameRequired, got %v", err)
	}
}

func TestUserUseCase_Create_NickNameTooLong(t *testing.T) {
	repo := &stubRepo{}
	uc := NewUserUseCase(repo)

	longNickName := strings.Repeat("a", 21)

	input := CreateUserInput{
		NickName:  longNickName,
		Email:     "test@example.com",
		Password:  "Password-hash123",
		IconImage: "https://example.com",
	}

	_, err := uc.Create(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != domain.ErrNickNameTooLong {
		t.Fatalf("expected ErrNickNameTooLong, got %v", err)
	}
}

func TestUserUseCase_Create_InvalidEmail(t *testing.T) {
	repo := &stubRepo{}
	uc := NewUserUseCase(repo)

	input := CreateUserInput{
		NickName:  "テストユーザー",
		Email:     "testexample.com",
		Password:  "Password-hash123",
		IconImage: "https://example.com",
	}

	_, err := uc.Create(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != domain.ErrInvalidEmail {
		t.Fatalf("expected ErrInvalidEmail, got %v", err)
	}
}

func TestUserUseCase_Create_PasswordTooShort(t *testing.T) {
	repo := &stubRepo{}
	uc := NewUserUseCase(repo)

	input := CreateUserInput{
		NickName:  "テストユーザー",
		Email:     "test@example.com",
		Password:  "p123@",
		IconImage: "https://example.com",
	}

	_, err := uc.Create(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != domain.ErrPasswordTooShort {
		t.Fatalf("expected ErrPasswordTooShort, got %v", err)
	}
}

func TestUserUseCase_Create_InvalidPassword(t *testing.T) {
	repo := &stubRepo{}
	uc := NewUserUseCase(repo)

	input := CreateUserInput{
		NickName:  "テストユーザー",
		Email:     "test@example.com",
		Password:  "Password",
		IconImage: "https://example.com",
	}

	_, err := uc.Create(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != domain.ErrInvalidPassword {
		t.Fatalf("expected ErrInvalidPassword, got %v", err)
	}
}
