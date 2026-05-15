package user

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aarondl/null/v8"
	"github.com/lib/pq"

	"notes-app/internal/apperror"
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

func TestMapInsertUserError(t *testing.T) {
	t.Parallel()

	plain := errors.New("plain db error")

	tests := []struct {
		name string
		err  error
		want error
	}{
		{
			name: "email unique violation",
			err: &pq.Error{
				Code:       pgSQLStateUniqueViolation,
				Constraint: constraintUsersEmailKey,
			},
			want: apperror.ErrEmailAlreadyExists,
		},
		{
			name: "nick_name unique violation",
			err: &pq.Error{
				Code:       pgSQLStateUniqueViolation,
				Constraint: constraintUsersNickNameKey,
			},
			want: apperror.ErrNickNameAlreadyTaken,
		},
		{
			name: "non pq error unchanged",
			err:  plain,
			want: plain,
		},
		{
			name: "nil unchanged",
			err:  nil,
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := mapInsertUserError(tt.err)
			if tt.want == nil {
				if got != nil {
					t.Fatalf("expected nil, got %v", got)
				}
				return
			}
			if !errors.Is(got, tt.want) {
				t.Fatalf("expected errors.Is(..., %v), got %v", tt.want, got)
			}
		})
	}

	t.Run("other unique violation returns original", func(t *testing.T) {
		t.Parallel()
		other := &pq.Error{
			Code:       pgSQLStateUniqueViolation,
			Constraint: "some_other_key",
		}
		got := mapInsertUserError(other)
		if got != other {
			t.Fatalf("expected same error back, got %v", got)
		}
	})
}

func TestUserRepository_Create_EmailDuplicate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewUser(db)
	fixture := models.User{
		ID:           "4c728e23-74d6-47b4-ae57-a9a3f2d242c7",
		NickName:     "テストユーザー",
		Email:        "dup@example.com",
		PasswordHash: "hash",
		IconImage:    null.StringFrom("https://example.com"),
	}

	mock.ExpectExec(`INSERT INTO "users"`).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnError(&pq.Error{
			Code:       pgSQLStateUniqueViolation,
			Constraint: constraintUsersEmailKey,
		})

	_, err = repo.Create(context.Background(), fixture)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, apperror.ErrEmailAlreadyExists) {
		t.Fatalf("expected ErrEmailAlreadyExists, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
