package apperror

var (
	// ErrOwnerNotFound は notes.user_id が存在しないユーザー（FK 違反）のときに返す。
	ErrOwnerNotFound = newAppError(
		CodeNoteOwnerNotFound,
		"ユーザーが見つかりません",
	)

	ErrTitleRequired = newAppError(
		CodeNoteTitleRequired,
		"タイトルを入力してください",
	)

	ErrTitleTooLong = newAppError(
		CodeNoteTitleTooLong,
		"タイトルは50文字以内で入力してください",
	)

	ErrContentRequired = newAppError(
		CodeNoteContentRequired,
		"本文を入力してください",
	)

	ErrContentTooLong = newAppError(
		CodeNoteContentTooLong,
		"本文は1000文字以内で入力してください",
	)
)
