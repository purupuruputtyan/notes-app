package note

import (
	"context"
	"errors"
	"strings"
	"testing"

	"notes-app/internal/models"
)

func TestNoteUseCase_Index_Success(t *testing.T) {
	userID := "11111111-1111-4111-8111-111111111111"
	repo := &StubRepo{
		notes: models.NoteSlice{
			{ID: "n1", UserID: userID, Title: "learn go", Content: "body1"},
			{ID: "n2", UserID: userID, Title: "second", Content: "body2"},
			{ID: "n3", UserID: "other-user", Title: "skip", Content: "x"},
		},
	}
	uc := NewNoteUseCase(repo)

	notes, err := uc.Index(context.Background(), userID)
	if err != nil {
		t.Fatalf("Index: %v", err)
	}
	if len(notes) != 2 {
		t.Fatalf("len: want 2, got %d", len(notes))
	}
	if notes[0].Title != "learn go" || notes[1].Title != "second" {
		t.Fatalf("titles: %+v / %+v", notes[0], notes[1])
	}
}

func TestNoteUseCase_Index_Empty(t *testing.T) {
	repo := &StubRepo{
		notes: models.NoteSlice{
			{ID: "n1", UserID: "aaaa", Title: "t", Content: "c"},
		},
	}
	uc := NewNoteUseCase(repo)

	notes, err := uc.Index(context.Background(), "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb")
	if err != nil {
		t.Fatalf("Index: %v", err)
	}
	if len(notes) != 0 {
		t.Fatalf("want no notes, got len=%d", len(notes))
	}
}

func TestNoteUseCase_Index_EmptyUserID(t *testing.T) {
	repo := &StubRepo{
		notes: models.NoteSlice{
			{ID: "n1", UserID: "", Title: "edge", Content: "c"},
		},
	}
	uc := NewNoteUseCase(repo)

	notes, err := uc.Index(context.Background(), "")
	if err != nil {
		t.Fatalf("Index: %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("len: want 1 for user_id empty match, got %d", len(notes))
	}
}

func TestNoteUseCase_Index_RepositoryError(t *testing.T) {
	dbErr := errors.New("db unavailable")
	repo := &StubRepo{indexErr: dbErr}
	uc := NewNoteUseCase(repo)

	_, err := uc.Index(context.Background(), "11111111-1111-4111-8111-111111111111")
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, dbErr) {
		t.Fatalf("want errors.Is(_, dbErr), got %v", err)
	}
	if !strings.Contains(err.Error(), "note usecase index") {
		t.Fatalf("want usecase wrap prefix in message, got %q", err.Error())
	}
}
