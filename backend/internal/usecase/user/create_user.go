package user

import (
	"github.com/aarondl/null/v8"
	"github.com/google/uuid"

	"notes-app/internal/models"
)

type CreateUserInput struct {
	NickName  string
	Email     string
	Password  string
	IconImage string
}

func (u *UserUseCase) Create(input CreateUserInput) (models.User, error) {
	if err := validateNickName(input.NickName); err != nil {
		return models.User{}, err
	}

	if err := validateEmail(input.Email); err != nil {
		return models.User{}, err
	}

	if err := validatePassword(input.Password); err != nil {
		return models.User{}, err
	}

	hashed, err := hashPassword(input.Password)

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
