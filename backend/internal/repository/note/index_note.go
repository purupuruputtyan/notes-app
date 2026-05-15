package note

import (
	"context"
	"fmt"

	"notes-app/internal/models"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

func (r *NoteRepository) Index(ctx context.Context, userID string) (models.NoteSlice, error) {
	rows, err := models.Notes(
		qm.OrderBy("created_at DESC"),
		models.NoteWhere.UserID.EQ(userID),
	).All(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("note repository index: %w", err)
	}

	return rows, nil
}
