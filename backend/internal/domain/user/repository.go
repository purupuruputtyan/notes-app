package user

import "notes-app/internal/models"

type Repository interface {
	Create(models.User) (models.User, error)
	Show(id string) (models.User, error)
}
