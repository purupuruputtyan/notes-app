package user

import (
	"context"
	"database/sql"
	"errors"

	"notes-app/internal/domain/user"
	"notes-app/internal/models"
)

func (r *UserRepository) Show(id string) (models.User, error) {
	ctx := context.Background()

	row, err := models.FindUser(ctx, r.db, id)
	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, user.ErrUserNotFound
	}
	if err != nil {
		return models.User{}, err
	}

	return *row, nil
}
