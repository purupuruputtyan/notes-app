package user

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aarondl/null/v8"
	"github.com/google/uuid"

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

	user := models.User{
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

	created, err := repo.Create(user)
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
			user.ID,
			user.NickName,
			user.Email,
			user.PasswordHash,
			user.IconImage,
			true,
		)

	mock.ExpectQuery(`(?i)select (.+) from "users"`).
		WithArgs(created.ID).
		WillReturnRows(rows)

	found, err := repo.Show(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if found.ID != user.ID {
		t.Fatalf("expected id to be set")
	}

	if found.NickName != user.NickName {
		t.Fatalf(
			"expected NickName %s, got %s",
			user.NickName,
			created.NickName,
		)
	}

	if found.Email != user.Email {
		t.Fatalf(
			"expected Email %s, got %s",
			user.Email,
			created.Email,
		)
	}

	if found.PasswordHash != user.PasswordHash {
		t.Fatalf(
			"expected PasswordHash %s, got %s",
			user.PasswordHash,
			created.PasswordHash,
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

	if _, err := uuid.Parse(created.ID); err != nil {
		t.Fatalf("invalid uuid: %s", created.ID)
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

	user := models.User{
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

	_, err = repo.Create(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	mock.ExpectQuery(`(?i)select (.+) from "users"`).
		WithArgs("not-found-id").
		WillReturnError(sql.ErrNoRows)

	_, err = repo.Show("not-found-id")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
