package note

import (
	"context"
	"strings"
	"testing"

	"notes-app/internal/apperror"
)

func TestNoteUseCase_Create(t *testing.T) {
	repo := &StubRepo{}
	uc := NewNoteUseCase(repo)
	userID := "11111111-1111-4111-8111-111111111111"

	input := CreateNoteInput{
		Title:   "テストタイトル",
		Content: "テスト本文",
	}

	created, err := uc.Create(context.Background(), userID, input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if created.ID == "" {
		t.Fatalf("expected ID to be set")
	}

	if created.UserID != userID {
		t.Fatalf("expected UserID %s, got %s", userID, created.UserID)
	}

	if created.Title != "テストタイトル" {
		t.Fatalf("expected Title テストタイトル, got %s", created.Title)
	}

	if created.Content != "テスト本文" {
		t.Fatalf("expected Content テスト本文, got %s", created.Content)
	}

	if len(repo.notes) != 1 {
		t.Fatalf("expected 1 note in repo, got %d", len(repo.notes))
	}
}

func TestNoteUseCase_Create_EmptyTitle(t *testing.T) {
	repo := &StubRepo{}
	uc := NewNoteUseCase(repo)

	input := CreateNoteInput{
		Title:   "",
		Content: "テスト本文",
	}

	_, err := uc.Create(context.Background(), "user-id", input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != apperror.ErrTitleRequired {
		t.Fatalf("expected ErrTitleRequired, got %v", err)
	}
}

func TestNoteUseCase_Create_WhitespaceOnlyTitle(t *testing.T) {
	repo := &StubRepo{}
	uc := NewNoteUseCase(repo)

	input := CreateNoteInput{
		Title:   "   ",
		Content: "テスト本文",
	}

	_, err := uc.Create(context.Background(), "user-id", input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != apperror.ErrTitleRequired {
		t.Fatalf("expected ErrTitleRequired, got %v", err)
	}
}

func TestNoteUseCase_Create_TitleTooLong(t *testing.T) {
	repo := &StubRepo{}
	uc := NewNoteUseCase(repo)

	longTitle := strings.Repeat("a", 51)

	input := CreateNoteInput{
		Title:   longTitle,
		Content: "テスト本文",
	}

	_, err := uc.Create(context.Background(), "user-id", input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != apperror.ErrTitleTooLong {
		t.Fatalf("expected ErrTitleTooLong, got %v", err)
	}
}

func TestNoteUseCase_Create_EmptyContent(t *testing.T) {
	repo := &StubRepo{}
	uc := NewNoteUseCase(repo)

	input := CreateNoteInput{
		Title:   "テストタイトル",
		Content: "",
	}

	_, err := uc.Create(context.Background(), "user-id", input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != apperror.ErrContentRequired {
		t.Fatalf("expected ErrContentRequired, got %v", err)
	}
}

func TestNoteUseCase_Create_ContentTooLong(t *testing.T) {
	repo := &StubRepo{}
	uc := NewNoteUseCase(repo)

	longContent := strings.Repeat("a", 1001)

	input := CreateNoteInput{
		Title:   "テストタイトル",
		Content: longContent,
	}

	_, err := uc.Create(context.Background(), "user-id", input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != apperror.ErrContentTooLong {
		t.Fatalf("expected ErrContentTooLong, got %v", err)
	}
}
