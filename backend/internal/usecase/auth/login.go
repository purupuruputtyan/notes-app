package auth

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"notes-app/internal/apperror"
	domain "notes-app/internal/domain/user"
	"notes-app/internal/lib/jwt"
)

type LoginUseCase struct {
	repo domain.Repository
}

func NewLoginUseCase(
	repo domain.Repository,
) *LoginUseCase {
	return &LoginUseCase{
		repo: repo,
	}
}

type LoginInput struct {
	Email    string
	Password string
}

func (u *LoginUseCase) Execute(
	ctx context.Context,
	input LoginInput,
) (string, error) {
	user, err := u.repo.FindByEmail(ctx, input.Email)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return "", apperror.ErrInvalidLogin
	}

	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
