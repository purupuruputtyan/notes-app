package user

import (
	"github.com/aarondl/null/v8"

	domain "notes-app/internal/domain/user"
	"notes-app/internal/models"
)

type UpdateUserInput struct {
	ID        string
	NickName  string
	Email     string
	Password  string
	IconImage string
}

func (u *UserUseCase) Update(input UpdateUserInput) (models.User, error) {
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

	params := domain.UpdateUserParams{
		NickName:     input.NickName,
		Email:        input.Email,
		PasswordHash: hashed,
		IconImage:    null.StringFrom(input.IconImage),
	}

	return u.repo.Update(input.ID, params)
}
