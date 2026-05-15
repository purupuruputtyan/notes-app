package note

import (
	"context"

	"github.com/aarondl/sqlboiler/v4/boil"

	"notes-app/internal/apperror"
	"notes-app/internal/models"
	"notes-app/internal/repository/pgerror"
)

func mapInsertNoteError(err error) error {
	pqErr, ok := pgerror.As(err)
	if !ok {
		return err
	}

	switch pqErr.Code {
	case pgerror.SQLStateForeignKeyViolation:
		if pqErr.Constraint == pgerror.ConstraintNotesUserIDFkey {
			return apperror.ErrOwnerNotFound
		}
	case pgerror.SQLStateUniqueViolation:
		if pqErr.Constraint == pgerror.ConstraintNotesPkey {
			return err
		}
	}

	return err
}

func (r *NoteRepository) Create(ctx context.Context, input models.Note) (models.Note, error) {
	row := &models.Note{
		ID:        input.ID,
		UserID:    input.UserID,
		Title:     input.Title,
		Content:   input.Content,
		CreatedAt: input.CreatedAt,
		UpdatedAt: input.UpdatedAt,
	}

	if err := row.Insert(ctx, r.db, boil.Infer()); err != nil {
		return models.Note{}, mapInsertNoteError(err)
	}

	return *row, nil
}
