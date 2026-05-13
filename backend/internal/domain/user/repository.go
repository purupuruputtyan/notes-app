package user

import (
	"notes-app/internal/models"

	"github.com/aarondl/null/v8"
)

type UpdateUserParams struct {
	ID           string
	NickName     string
	Email        string
	PasswordHash string
	IconImage    null.String
}

type Repository interface {
	Create(models.User) (models.User, error)
	Show(id string) (models.User, error)
	Update(id string, params UpdateUserParams) (models.User, error)
}
