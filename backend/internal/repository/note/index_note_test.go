package note

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNoteRepository_Index_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	repo := NewNote(db)
	userID := "11111111-1111-4111-8111-111111111111"
	ts := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "title", "content", "created_at", "updated_at",
	}).
		AddRow("aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa", userID, "t1", "c1", ts, ts).
		AddRow("bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb", userID, "t2", "c2", ts, ts)

	mock.ExpectQuery(`(?i)select (.+) from "notes"`).
		WithArgs(userID).
		WillReturnRows(rows)

	got, err := repo.Index(context.Background(), userID)
	if err != nil {
		t.Fatalf("Index: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len: want 2, got %d", len(got))
	}
	if got[0].ID != "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa" {
		t.Fatalf("first id: %q", got[0].ID)
	}
	if got[0].UserID != userID || got[0].Title != "t1" || got[0].Content != "c1" {
		t.Fatalf("first row: %+v", got[0])
	}
	if got[1].UserID != userID || got[1].Title != "t2" || got[1].Content != "c2" {
		t.Fatalf("second row: %+v", got[1])
	}
	if !got[0].CreatedAt.Equal(ts) || !got[1].CreatedAt.Equal(ts) {
		t.Fatalf("CreatedAt: want %v, got %v / %v", ts, got[0].CreatedAt, got[1].CreatedAt)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
}

func TestNoteRepository_Index_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	repo := NewNote(db)
	userID := "22222222-2222-4222-8222-222222222222"

	mock.ExpectQuery(`(?i)select (.+) from "notes"`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "title", "content", "created_at", "updated_at",
		}))

	got, err := repo.Index(context.Background(), userID)
	if err != nil {
		t.Fatalf("Index: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("want empty slice, got len=%d", len(got))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
}

func TestNoteRepository_Index_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	repo := NewNote(db)
	userID := "33333333-3333-4333-8333-333333333333"
	dbErr := errors.New("connection reset")

	mock.ExpectQuery(`(?i)select (.+) from "notes"`).
		WithArgs(userID).
		WillReturnError(dbErr)

	_, err = repo.Index(context.Background(), userID)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, dbErr) {
		t.Fatalf("Index: want errors.Is(_, dbErr), got %v", err)
	}
	if !strings.Contains(err.Error(), "note repository index") {
		t.Fatalf("want repository wrap prefix in message, got %q", err.Error())
	}
	if mockErr := mock.ExpectationsWereMet(); mockErr != nil {
		t.Fatalf("sqlmock: %v", mockErr)
	}
}
