package user

import (
	"context"

	"github.com/aarondl/sqlboiler/v4/boil"

	"notes-app/internal/models"
)

func (r *UserRepository) Create(input models.User) (models.User, error) {
	ctx := context.Background()

	row := &models.User{
		ID:           input.ID,
		NickName:     input.NickName,
		Email:        input.Email,
		PasswordHash: input.PasswordHash,
		IconImage:    input.IconImage,
		IsActive:     true,
	}

	if err := row.Insert(ctx, r.db, boil.Infer()); err != nil {
		return models.User{}, err
	}

	return *row, nil
}
