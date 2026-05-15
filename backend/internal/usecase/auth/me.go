package auth

import (
	"context"

	domain "notes-app/internal/domain/user"
	"notes-app/internal/models"
)

type MeUseCase struct {
	repo domain.Repository
}

func NewMeUseCase(
	repo domain.Repository,
) *MeUseCase {
	return &MeUseCase{
		repo: repo,
	}
}

func (u *MeUseCase) Execute(
	ctx context.Context,
	userID string,
) (models.User, error) {
	return u.repo.Show(ctx, userID)
}
