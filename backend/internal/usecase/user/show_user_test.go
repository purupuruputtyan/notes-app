package user

import (
	"errors"
	"testing"

	domain "notes-app/internal/domain/user"
)

func TestUserUseCase_Show(t *testing.T) {
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
		t.Fatalf("Create: %v", err)
	}

	found, err := uc.Show(created.ID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if found.ID != created.ID {
		t.Fatalf(
			"expected ID %s, got %s",
			created.ID,
			found.ID,
		)
	}

	if found.NickName != created.NickName {
		t.Fatalf(
			"expected NickName %s, got %s",
			created.NickName,
			found.NickName,
		)
	}

	if found.Email != created.Email {
		t.Fatalf(
			"expected Email %s, got %s",
			created.Email,
			found.Email,
		)
	}

	if found.PasswordHash != created.PasswordHash {
		t.Fatalf(
			"expected PasswordHash %s, got %s",
			created.PasswordHash,
			found.PasswordHash,
		)
	}

	if found.IconImage != created.IconImage {
		t.Fatalf(
			"expected IconImage %v, got %v",
			created.IconImage,
			found.IconImage,
		)
	}

	if found.IsActive != created.IsActive {
		t.Fatalf(
			"expected IsActive %v, got %v",
			created.IsActive,
			found.IsActive,
		)
	}
}

func TestUserUseCase_Show_NotFound(t *testing.T) {
	repo := &stubRepo{}
	uc := NewUserUseCase(repo)

	_, err := uc.Show("not-found-id")

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
