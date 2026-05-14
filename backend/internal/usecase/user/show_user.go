package user

import (
	"context"

	"notes-app/internal/models"
)

func (u *UserUseCase) Show(ctx context.Context, id string) (models.User, error) {
	return u.repo.Show(ctx, id)
}
