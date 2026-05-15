package apperror

var (
	// ErrOwnerNotFound は notes.user_id が存在しないユーザー（FK 違反）のときに返す。
	ErrOwnerNotFound = newAppError(
		CodeNoteOwnerNotFound,
		"ユーザーが見つかりません",
	)
)
