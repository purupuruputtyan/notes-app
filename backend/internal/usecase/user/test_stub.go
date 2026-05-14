package user

import (
	"context"
	domain "notes-app/internal/domain/user"
	"notes-app/internal/models"
)

// StubRepo はユースケース／ハンドラのテスト用インメモリ実装（本番コードからは使わないこと）。
type StubRepo struct {
	users []models.User
}

func (s *StubRepo) Create(ctx context.Context, u models.User) (models.User, error) {
	s.users = append(s.users, u)
	return u, nil
}

func (s *StubRepo) Show(ctx context.Context, id string) (models.User, error) {
	for _, u := range s.users {
		if u.ID == id {
			return u, nil
		}
	}

	return models.User{}, domain.ErrUserNotFound
}

func (s *StubRepo) Update(
	ctx context.Context,
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

func (s *StubRepo) FindByEmail(ctx context.Context, email string) (models.User, error) {
	for _, u := range s.users {
		if u.Email == email {
			return u, nil
		}
	}

	return models.User{}, domain.ErrUserNotFound
}
