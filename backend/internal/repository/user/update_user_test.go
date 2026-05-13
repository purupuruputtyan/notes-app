package user

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aarondl/null/v8"

	domain "notes-app/internal/domain/user"
	"notes-app/internal/models"
)

func TestUserRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewUser(db)

	user := models.User{
		ID:           "4c728e23-74d6-47b4-ae57-a9a3f2d242c7",
		NickName:     "テストユーザー",
		Email:        "test@example.com",
		PasswordHash: "password-hash",
		IconImage:    null.StringFrom("https://example.com"),
		IsActive:     true,
	}

	// FindUser 用
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
			user.IsActive,
		)

	mock.ExpectQuery(`select \* from "users" where "id"=\$1`).
		WithArgs(user.ID).
		WillReturnRows(rows)

	// UPDATE 用
	mock.ExpectExec(`UPDATE "users"`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	updated, err := repo.Update(UpdateUserInput{
		ID:           user.ID,
		NickName:     user.NickName,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		IconImage:    user.IconImage,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updated.NickName != user.NickName {
		t.Fatalf("expected NickName %s, got %s", user.NickName, updated.NickName)
	}

	if updated.Email != user.Email {
		t.Fatalf("expected Email %s, got %s", user.Email, updated.Email)
	}

	if updated.PasswordHash != user.PasswordHash {
		t.Fatalf("expected PasswordHash %s, got %s", user.PasswordHash, updated.PasswordHash)
	}

	if updated.IconImage != user.IconImage {
		t.Fatalf("expected IconImage %v, got %v", user.IconImage, updated.IconImage)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func TestUserRepository_Update_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewUser(db)

	mock.ExpectQuery(`select \* from "users" where "id"=\$1`).
		WithArgs("not-found-id").
		WillReturnError(sql.ErrNoRows)

	_, err = repo.Update(UpdateUserInput{
		ID: "not-found-id",
	})

	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
