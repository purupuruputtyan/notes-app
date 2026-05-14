package user

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aarondl/null/v8"

	domain "notes-app/internal/domain/user"
	"notes-app/internal/models"
)

func TestUserRepository_Show(
	t *testing.T,
) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(
			"failed to create sqlmock: %v",
			err,
		)
	}

	defer db.Close()

	repo := NewUser(db)

	fixture := models.User{
		ID:           "4c728e23-74d6-47b4-ae57-a9a3f2d242c7",
		NickName:     "テストユーザー",
		Email:        "test@example.com",
		PasswordHash: "password-hash",
		IconImage:    null.StringFrom("https://example.com"),
	}

	mock.ExpectExec(`INSERT INTO "users"`).
		WithArgs(
			sqlmock.AnyArg(), // ID
			sqlmock.AnyArg(), // NickName
			sqlmock.AnyArg(), // Email
			sqlmock.AnyArg(), // PasswordHash
			sqlmock.AnyArg(), // IconImage
			sqlmock.AnyArg(), // IsActive
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	created, err := repo.Create(context.Background(), fixture)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	rows := sqlmock.NewRows([]string{
		"id",
		"nick_name",
		"email",
		"password_hash",
		"icon_image",
		"is_active",
	}).
		AddRow(
			fixture.ID,
			fixture.NickName,
			fixture.Email,
			fixture.PasswordHash,
			fixture.IconImage,
			true,
		)

	mock.ExpectQuery(`(?i)select (.+) from "users"`).
		WithArgs(created.ID).
		WillReturnRows(rows)

	found, err := repo.Show(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if found.ID != fixture.ID {
		t.Fatalf(
			"expected ID %s, got %s",
			fixture.ID,
			found.ID,
		)
	}

	if found.NickName != fixture.NickName {
		t.Fatalf(
			"expected NickName %s, got %s",
			fixture.NickName,
			found.NickName,
		)
	}

	if found.Email != fixture.Email {
		t.Fatalf(
			"expected Email %s, got %s",
			fixture.Email,
			found.Email,
		)
	}

	if found.PasswordHash != fixture.PasswordHash {
		t.Fatalf(
			"expected PasswordHash %s, got %s",
			fixture.PasswordHash,
			found.PasswordHash,
		)
	}

	if !found.IsActive {
		t.Fatalf("expected IsActive to be true")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(
			"unmet sql expectations: %v",
			err,
		)
	}
}

func TestUserRepository_Show_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(
			"failed to create sqlmock: %v",
			err,
		)
	}

	defer db.Close()

	repo := NewUser(db)

	fixture := models.User{
		ID:           "4c728e23-74d6-47b4-ae57-a9a3f2d242c7",
		NickName:     "テストユーザー",
		Email:        "test@example.com",
		PasswordHash: "password-hash",
		IconImage:    null.StringFrom("https://example.com"),
	}

	mock.ExpectExec(`INSERT INTO "users"`).
		WithArgs(
			sqlmock.AnyArg(), // ID
			sqlmock.AnyArg(), // NickName
			sqlmock.AnyArg(), // Email
			sqlmock.AnyArg(), // PasswordHash
			sqlmock.AnyArg(), // IconImage
			sqlmock.AnyArg(), // IsActive
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = repo.Create(context.Background(), fixture)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	mock.ExpectQuery(`(?i)select (.+) from "users"`).
		WithArgs("not-found-id").
		WillReturnError(sql.ErrNoRows)

	_, err = repo.Show(context.Background(), "not-found-id")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
