package note

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"

	"notes-app/internal/apperror"
	"notes-app/internal/models"
	"notes-app/internal/repository/pgerror"
)

func TestNoteRepository_Create_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	repo := NewNote(db)
	ts := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	fixture := models.Note{
		ID:        "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa",
		UserID:    "11111111-1111-4111-8111-111111111111",
		Title:     "title",
		Content:   "content",
		CreatedAt: ts,
		UpdatedAt: ts,
	}

	mock.ExpectExec(`INSERT INTO "notes"`).
		WithArgs(
			fixture.ID,
			fixture.UserID,
			fixture.Title,
			fixture.Content,
			fixture.CreatedAt,
			fixture.UpdatedAt,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	created, err := repo.Create(context.Background(), fixture)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.ID != fixture.ID {
		t.Fatalf("ID: want %q, got %q", fixture.ID, created.ID)
	}
	if created.UserID != fixture.UserID || created.Title != fixture.Title || created.Content != fixture.Content {
		t.Fatalf("row: %+v", created)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestMapInsertNoteError(t *testing.T) {
	t.Parallel()

	plain := errors.New("plain db error")

	tests := []struct {
		name string
		err  error
		want error
	}{
		{
			name: "owner user_id foreign key violation",
			err: &pq.Error{
				Code:       pgerror.SQLStateForeignKeyViolation,
				Constraint: pgerror.ConstraintNotesUserIDFkey,
			},
			want: apperror.ErrOwnerNotFound,
		},
		{
			name: "notes_pkey unique violation returns original",
			err: &pq.Error{
				Code:       pgerror.SQLStateUniqueViolation,
				Constraint: pgerror.ConstraintNotesPkey,
			},
			want: &pq.Error{
				Code:       pgerror.SQLStateUniqueViolation,
				Constraint: pgerror.ConstraintNotesPkey,
			},
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
			got := mapInsertNoteError(tt.err)
			if tt.want == nil {
				if got != nil {
					t.Fatalf("expected nil, got %v", got)
				}
				return
			}
			if pqWant, ok := tt.want.(*pq.Error); ok {
				pqGot, ok := got.(*pq.Error)
				if !ok || pqGot.Code != pqWant.Code || pqGot.Constraint != pqWant.Constraint {
					t.Fatalf("expected pq.Error %+v, got %v", pqWant, got)
				}
				return
			}
			if !errors.Is(got, tt.want) {
				t.Fatalf("expected errors.Is(..., %v), got %v", tt.want, got)
			}
		})
	}

	t.Run("other foreign key violation returns original", func(t *testing.T) {
		t.Parallel()
		other := &pq.Error{
			Code:       pgerror.SQLStateForeignKeyViolation,
			Constraint: "some_other_fkey",
		}
		got := mapInsertNoteError(other)
		if got != other {
			t.Fatalf("expected same error back, got %v", got)
		}
	})

	t.Run("other unique violation returns original", func(t *testing.T) {
		t.Parallel()
		other := &pq.Error{
			Code:       pgerror.SQLStateUniqueViolation,
			Constraint: "some_other_key",
		}
		got := mapInsertNoteError(other)
		if got != other {
			t.Fatalf("expected same error back, got %v", got)
		}
	})
}

func TestNoteRepository_Create_OwnerNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	repo := NewNote(db)
	ts := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	fixture := models.Note{
		ID:        "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa",
		UserID:    "00000000-0000-4000-8000-000000000000",
		Title:     "title",
		Content:   "content",
		CreatedAt: ts,
		UpdatedAt: ts,
	}

	mock.ExpectExec(`INSERT INTO "notes"`).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnError(&pq.Error{
			Code:       pgerror.SQLStateForeignKeyViolation,
			Constraint: pgerror.ConstraintNotesUserIDFkey,
		})

	_, err = repo.Create(context.Background(), fixture)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, apperror.ErrOwnerNotFound) {
		t.Fatalf("expected ErrOwnerNotFound, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
