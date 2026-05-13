package user

import (
	"notes-app/internal/models"
)

func (u *UserUseCase) Show(id string) (models.User, error) {
	return u.repo.Show(id)
}
