package user

import (
	"context"
	"database/sql"
	"errors"

	"notes-app/internal/domain/user"
	"notes-app/internal/models"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
)

type UpdateUserInput struct {
	ID           string
	NickName     string
	Email        string
	PasswordHash string
	IconImage    null.String
}

func (r *UserRepository) Update(input UpdateUserInput) (models.User, error) {
	ctx := context.Background()

	row, err := models.FindUser(ctx, r.db, input.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, user.ErrUserNotFound
	}
	if err != nil {
		return models.User{}, err
	}

	row.NickName = input.NickName
	row.Email = input.Email
	row.PasswordHash = input.PasswordHash
	row.IconImage = input.IconImage

	rowsAffected, err := row.Update(
		ctx,
		r.db,
		boil.Whitelist(
			models.UserColumns.NickName,
			models.UserColumns.Email,
			models.UserColumns.PasswordHash,
			models.UserColumns.IconImage,
		),
	)

	if err != nil {
		return models.User{}, err
	}
	if rowsAffected == 0 {
		return models.User{}, user.ErrUserNotFound
	}

	return *row, nil
}
