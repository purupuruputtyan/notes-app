package note

import (
	"context"
	"fmt"

	"notes-app/internal/models"
)

// Index はログインユーザーに紐づくノート一覧を返す。
// userID はハンドラで認可ミドルウェアから取得した主体 ID を渡すこと（リクエストパラメータの任意値をそのまま渡さない）。
func (u *NoteUseCase) Index(ctx context.Context, userID string) (models.NoteSlice, error) {
	notes, err := u.repo.Index(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("note usecase index: %w", err)
	}

	return notes, nil
}
