package note

import (
	"database/sql"
)

type NoteRepository struct {
	db *sql.DB
}

func NewNote(db *sql.DB) *NoteRepository {
	return &NoteRepository{db: db}
}
