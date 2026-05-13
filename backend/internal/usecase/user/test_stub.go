package user

import (
	domain "notes-app/internal/domain/user"
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
	for _, u := range s.users {
		if u.ID == id {
			return u, nil
		}
	}

	return models.User{}, domain.ErrUserNotFound
}

func (s *stubRepo) Update(
	id string,
	u domain.UpdateUserParams,
) (models.User, error) {

	for i, existingUser := range s.users {

		if existingUser.ID == id {

			s.users[i].NickName = u.NickName
			s.users[i].Email = u.Email
			s.users[i].PasswordHash = u.PasswordHash
			s.users[i].IconImage = u.IconImage

			return s.users[i], nil
		}
	}

	return models.User{}, domain.ErrUserNotFound
}
