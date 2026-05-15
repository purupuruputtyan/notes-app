package note

import (
	"context"

	"notes-app/internal/models"
)

type Repository interface {
	// Index は user_id が一致するノート一覧を返す。
	// userID には認証済み主体の ID（例: JWT の sub / user_id）を渡すこと。クライアントが任意に指定できる値だけを信頼しないこと。
	Index(ctx context.Context, userID string) (models.NoteSlice, error)
	Create(ctx context.Context, input models.Note) (models.Note, error)
}
