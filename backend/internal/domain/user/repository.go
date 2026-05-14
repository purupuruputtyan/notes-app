package user

import (
	"context"

	"notes-app/internal/models"

	"github.com/aarondl/null/v8"
)

type UpdateUserParams struct {
	NickName     string
	Email        string
	PasswordHash string
	IconImage    null.String
}

type Repository interface {
	Create(ctx context.Context, input models.User) (models.User, error)
	Show(ctx context.Context, id string) (models.User, error)
	Update(ctx context.Context, id string, params UpdateUserParams) (models.User, error)
	FindByEmail(ctx context.Context, email string) (models.User, error)
}
