package user

import (
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

	user := models.User{
		NickName:     "テストユーザー",
		Email:        "test@example.com",
		PasswordHash: "password-hash",
		IconImage:    null.StringFrom("https://example.com"),
	}

	mock.ExpectExec(`INSERT INTO "users"`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	created, err := repo.Create(user)

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	if created.ID == "" {
		t.Fatalf("expected id to be set")
	}

	if created.NickName != user.NickName {
		t.Fatalf(
			"expected NickName %s, got %s",
			user.NickName,
			created.NickName,
		)
	}

	if created.Email != user.Email {
		t.Fatalf(
			"expected Email %s, got %s",
			user.Email,
			created.Email,
		)
	}

	if created.PasswordHash != user.PasswordHash {
		t.Fatalf(
			"expected PasswordHash %s, got %s",
			user.PasswordHash,
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
