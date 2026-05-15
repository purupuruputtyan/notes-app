package user

import (
	"context"

	"github.com/aarondl/sqlboiler/v4/boil"

	"notes-app/internal/apperror"
	"notes-app/internal/models"
	"notes-app/internal/repository/pgerror"
)

func mapInsertUserError(err error) error {
	pqErr, ok := pgerror.As(err)
	if !ok {
		return err
	}
	if pqErr.Code != pgerror.SQLStateUniqueViolation {
		return err
	}

	switch pqErr.Constraint {
	case pgerror.ConstraintUsersEmailKey:
		return apperror.ErrEmailAlreadyExists
	case pgerror.ConstraintUsersNickNameKey:
		return apperror.ErrNickNameAlreadyTaken
	default:
		return err
	}
}

func (r *UserRepository) Create(ctx context.Context, input models.User) (models.User, error) {
	row := &models.User{
		ID:           input.ID,
		NickName:     input.NickName,
		Email:        input.Email,
		PasswordHash: input.PasswordHash,
		IconImage:    input.IconImage,
		IsActive:     true,
	}

	if err := row.Insert(ctx, r.db, boil.Infer()); err != nil {
		return models.User{}, mapInsertUserError(err)
	}

	return *row, nil
}
