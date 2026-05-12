package user

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/aarondl/null/v8"
	"github.com/google/uuid"

	"notes-app/internal/domain/user"
	"notes-app/internal/models"
)

type UserUseCase struct {
	repo user.Repository
}

type CreateUserInput struct {
	NickName     string
	Email        string
	PasswordHash string
	IconImage    string
}

func NewUserUseCase(repo user.Repository) *UserUseCase {
	return &UserUseCase{
		repo: repo,
	}
}

func (u *UserUseCase) Create(input CreateUserInput) (models.User, error) {
	if err := validateNickName(input.NickName); err != nil {
		return models.User{}, err
	}

	if err := validateEmail(input.Email); err != nil {
		return models.User{}, err
	}

	if err := validatePassword(input.PasswordHash); err != nil {
		return models.User{}, err
	}

	hashed, err := hashPassword(input.PasswordHash)

	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		ID:           uuid.NewString(),
		NickName:     input.NickName,
		Email:        input.Email,
		PasswordHash: hashed,
		IconImage:    null.StringFrom(input.IconImage),
	}
	return u.repo.Create(user)
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
