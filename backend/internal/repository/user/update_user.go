package user

import (
	"context"
	"database/sql"
	"errors"

	"notes-app/internal/domain/user"
	"notes-app/internal/models"

	"github.com/aarondl/sqlboiler/v4/boil"
)

func (r *UserRepository) Update(ctx context.Context, id string, input user.UpdateUserParams) (models.User, error) {
	row, err := models.FindUser(ctx, r.db, id)
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

	_, err = row.Update(
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

	return *row, nil
}
