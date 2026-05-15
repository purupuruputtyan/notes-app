package user

import (
	"context"
	"database/sql"
	"errors"

	"notes-app/internal/apperror"
	"notes-app/internal/models"
)

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (models.User, error) {
	row, err := models.Users(models.UserWhere.Email.EQ(email)).One(ctx, r.db)
	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, apperror.ErrUserNotFound
	}
	if err != nil {
		return models.User{}, err
	}

	return *row, nil
}
