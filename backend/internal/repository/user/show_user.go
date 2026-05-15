package user

import (
	"context"
	"database/sql"
	"errors"

	"notes-app/internal/apperror"
	"notes-app/internal/models"
)

func (r *UserRepository) Show(ctx context.Context, id string) (models.User, error) {
	row, err := models.FindUser(ctx, r.db, id)
	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, apperror.ErrUserNotFound
	}
	if err != nil {
		return models.User{}, err
	}

	return *row, nil
}
