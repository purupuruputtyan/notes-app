package user

import (
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUser(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}
