package note

import (
	"notes-app/internal/domain/note"
)

type NoteUseCase struct {
	repo note.Repository
}

func NewNoteUseCase(repo note.Repository) *NoteUseCase {
	return &NoteUseCase{
		repo: repo,
	}
}
