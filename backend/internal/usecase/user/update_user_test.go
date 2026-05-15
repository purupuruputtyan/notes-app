package user

import (
	"context"
	"strings"
	"testing"

	"notes-app/internal/apperror"
)

func TestUserUseCase_Update(t *testing.T) {
	repo := &StubRepo{}
	uc := NewUserUseCase(repo)

	created, err := uc.Create(context.Background(), CreateUserInput{
		NickName:  "テストユーザー",
		Email:     "test@example.com",
		Password:  "Password123!",
		IconImage: "https://example.com",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	input := UpdateUserInput{
		ID:        created.ID,
		NickName:  "アップデートユーザー",
		Email:     "update@example.com",
		Password:  "update-hash123",
		IconImage: "https://update.com",
	}

	updated, err := uc.Update(context.Background(), input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updated.ID != input.ID {
		t.Fatalf("expected ID, got %s", updated.ID)
	}

	if updated.NickName != input.NickName {
		t.Fatalf("expected NickName アップデートユーザー, got %s", updated.NickName)
	}

	if updated.Email != input.Email {
		t.Fatalf("expected Email update@example.com, got %s", updated.Email)
	}

	if updated.PasswordHash == "" {
		t.Fatalf("expected Password to be set")
	}

	if updated.PasswordHash == created.PasswordHash {
		t.Fatalf("Password must be hashed")
	}

	if updated.IconImage.String != input.IconImage {
		t.Fatalf("expected IconImage https://update.com, got %s", input.IconImage)
	}

	if len(repo.users) != 1 {
		t.Fatalf("expected 1 user in repo, got %d", len(repo.users))
	}
}

func TestUserUseCase_Update_EmptyNickName(t *testing.T) {
	repo := &StubRepo{}
	uc := NewUserUseCase(repo)

	created, err := uc.Create(context.Background(), CreateUserInput{
		NickName:  "テストユーザー",
		Email:     "test@example.com",
		Password:  "Password123!",
		IconImage: "https://example.com",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	input := UpdateUserInput{
		ID:        created.ID,
		NickName:  "",
		Email:     "test@example.com",
		Password:  "Password-hash123",
		IconImage: "https://example.com",
	}

	_, err = uc.Update(context.Background(), input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != apperror.ErrNickNameRequired {
		t.Fatalf("expected ErrNickNameRequired, got %v", err)
	}
}

func TestUserUseCase_Update_NickNameTooLong(t *testing.T) {
	repo := &StubRepo{}
	uc := NewUserUseCase(repo)

	created, err := uc.Create(context.Background(), CreateUserInput{
		NickName:  "テストユーザー",
		Email:     "test@example.com",
		Password:  "Password123!",
		IconImage: "https://example.com",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	longNickName := strings.Repeat("a", 21)

	input := UpdateUserInput{
		ID:        created.ID,
		NickName:  longNickName,
		Email:     "test@example.com",
		Password:  "Password-hash123",
		IconImage: "https://example.com",
	}

	_, err = uc.Update(context.Background(), input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != apperror.ErrNickNameTooLong {
		t.Fatalf("expected ErrNickNameTooLong, got %v", err)
	}
}

func TestUserUseCase_Update_InvalidEmail(t *testing.T) {
	repo := &StubRepo{}
	uc := NewUserUseCase(repo)

	created, err := uc.Create(context.Background(), CreateUserInput{
		NickName:  "テストユーザー",
		Email:     "test@example.com",
		Password:  "Password123!",
		IconImage: "https://example.com",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	input := UpdateUserInput{
		ID:        created.ID,
		NickName:  "テストユーザー",
		Email:     "testexample.com",
		Password:  "Password-hash123",
		IconImage: "https://example.com",
	}

	_, err = uc.Update(context.Background(), input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != apperror.ErrInvalidEmail {
		t.Fatalf("expected ErrInvalidEmail, got %v", err)
	}
}

func TestUserUseCase_Update_PasswordTooShort(t *testing.T) {
	repo := &StubRepo{}
	uc := NewUserUseCase(repo)

	created, err := uc.Create(context.Background(), CreateUserInput{
		NickName:  "テストユーザー",
		Email:     "test@example.com",
		Password:  "Password123!",
		IconImage: "https://example.com",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	input := UpdateUserInput{
		ID:        created.ID,
		NickName:  "テストユーザー",
		Email:     "test@example.com",
		Password:  "p123@",
		IconImage: "https://example.com",
	}

	_, err = uc.Update(context.Background(), input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != apperror.ErrPasswordTooShort {
		t.Fatalf("expected ErrPasswordTooShort, got %v", err)
	}
}

func TestUserUseCase_Update_InvalidPassword(t *testing.T) {
	repo := &StubRepo{}
	uc := NewUserUseCase(repo)

	created, err := uc.Create(context.Background(), CreateUserInput{
		NickName:  "テストユーザー",
		Email:     "test@example.com",
		Password:  "Password123!",
		IconImage: "https://example.com",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	input := UpdateUserInput{
		ID:        created.ID,
		NickName:  "テストユーザー",
		Email:     "test@example.com",
		Password:  "Password",
		IconImage: "https://example.com",
	}

	_, err = uc.Update(context.Background(), input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != apperror.ErrInvalidPassword {
		t.Fatalf("expected ErrInvalidPassword, got %v", err)
	}
}
