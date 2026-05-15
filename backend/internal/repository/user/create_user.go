package user

import (
	"context"
	"errors"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/lib/pq"

	"notes-app/internal/apperror"
	"notes-app/internal/models"
)

func mapInsertUserError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == pgSQLStateUniqueViolation {
		switch pqErr.Constraint {
		case constraintUsersEmailKey:
			return apperror.ErrEmailAlreadyExists
		case constraintUsersNickNameKey:
			return apperror.ErrNickNameAlreadyTaken
		default:
			return err
		}
	}
	return err
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
