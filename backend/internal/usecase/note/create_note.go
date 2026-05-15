package note

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"notes-app/internal/models"
)

type CreateNoteInput struct {
	Title   string
	Content string
}

// Create は認証済みユーザー userID に紐づくノートを作成する。
// userID はハンドラで認可ミドルウェアから取得した主体 ID を渡すこと。
func (u *NoteUseCase) Create(ctx context.Context, userID string, input CreateNoteInput) (models.Note, error) {
	if err := validateTitle(input.Title); err != nil {
		return models.Note{}, err
	}

	if err := validateContent(input.Content); err != nil {
		return models.Note{}, err
	}

	title := strings.TrimSpace(input.Title)
	content := strings.TrimSpace(input.Content)

	note := models.Note{
		ID:      uuid.NewString(),
		UserID:  userID,
		Title:   title,
		Content: content,
	}
	return u.repo.Create(ctx, note)
}
