package user

import (
	"golang.org/x/crypto/bcrypt"

	"notes-app/internal/domain/user"
)

type UserUseCase struct {
	repo user.Repository
}

func NewUserUseCase(repo user.Repository) *UserUseCase {
	return &UserUseCase{
		repo: repo,
	}
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}
