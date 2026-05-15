package user

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aarondl/null/v8"

	"notes-app/internal/apperror"
	"notes-app/internal/models"
)

// sqlboiler が発行する「email で 1 件」の SELECT（id 検索と取り違えないようにする）。
const findByEmailQueryPattern = `(?i)^SELECT "users"\.\* FROM "users" WHERE \("users"\."email"\s*=\s*\$1\) LIMIT 1;$`

func TestUserRepository_FindByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewUser(db)

	email := "test@example.com"
	fixture := models.User{
		ID:           "4c728e23-74d6-47b4-ae57-a9a3f2d242c7",
		NickName:     "テストユーザー",
		Email:        email,
		PasswordHash: "password-hash",
		IconImage:    null.StringFrom("https://example.com"),
		IsActive:     true,
		CreatedAt:    time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC),
		UpdatedAt:    time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC),
	}

	rows := sqlmock.NewRows([]string{
		"id",
		"nick_name",
		"email",
		"password_hash",
		"icon_image",
		"is_active",
		"created_at",
		"updated_at",
	}).
		AddRow(
			fixture.ID,
			fixture.NickName,
			fixture.Email,
			fixture.PasswordHash,
			fixture.IconImage,
			fixture.IsActive,
			fixture.CreatedAt,
			fixture.UpdatedAt,
		)

	mock.ExpectQuery(findByEmailQueryPattern).
		WithArgs(email).
		WillReturnRows(rows)

	found, err := repo.FindByEmail(context.Background(), email)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if found.ID != fixture.ID {
		t.Fatalf("expected ID %s, got %s", fixture.ID, found.ID)
	}
	if found.NickName != fixture.NickName {
		t.Fatalf("expected NickName %s, got %s", fixture.NickName, found.NickName)
	}
	if found.Email != fixture.Email {
		t.Fatalf("expected Email %s, got %s", fixture.Email, found.Email)
	}
	if found.PasswordHash != fixture.PasswordHash {
		t.Fatalf("expected PasswordHash %s, got %s", fixture.PasswordHash, found.PasswordHash)
	}
	if !found.IsActive {
		t.Fatalf("expected IsActive to be true")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewUser(db)

	email := "missing@example.com"

	mock.ExpectQuery(findByEmailQueryPattern).
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)

	_, err = repo.FindByEmail(context.Background(), email)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, apperror.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}
