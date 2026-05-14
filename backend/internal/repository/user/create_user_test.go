package user

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aarondl/null/v8"

	"notes-app/internal/models"
)

func TestUserRepository_Create(
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
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	if created.ID == "" {
		t.Fatalf("expected id to be set")
	}

	if created.NickName != fixture.NickName {
		t.Fatalf(
			"expected NickName %s, got %s",
			fixture.NickName,
			created.NickName,
		)
	}

	if created.Email != fixture.Email {
		t.Fatalf(
			"expected Email %s, got %s",
			fixture.Email,
			created.Email,
		)
	}

	if created.PasswordHash != fixture.PasswordHash {
		t.Fatalf(
			"expected PasswordHash %s, got %s",
			fixture.PasswordHash,
			created.PasswordHash,
		)
	}

	if !created.IsActive {
		t.Fatalf("expected IsActive to be true")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(
			"unmet sql expectations: %v",
			err,
		)
	}
}
