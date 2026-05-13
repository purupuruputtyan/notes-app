package user

import (
	"notes-app/internal/domain/user"
	"notes-app/internal/models"
)

type stubRepo struct {
	users []models.User
}

func (s *stubRepo) Create(u models.User) (models.User, error) {
	s.users = append(s.users, u)
	return u, nil
}

func (s *stubRepo) Show(id string) (models.User, error) {
	for _, user := range s.users {
		if user.ID == id {
			return user, nil
		}
	}

	return models.User{}, user.ErrUserNotFound
}
