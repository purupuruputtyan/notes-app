package note

import (
	"context"

	"notes-app/internal/models"
)

type Repository interface {
	Index(ctx context.Context, userID string) (models.NoteSlice, error)
}
