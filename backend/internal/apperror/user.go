package apperror

var (
	ErrNickNameRequired = newAppError(
		CodeUserNickNameRequired,
		"ニックネームを入力してください",
	)

	ErrNickNameTooLong = newAppError(
		CodeUserNickNameTooLong,
		"ニックネームは20文字以内で入力してください",
	)

	ErrInvalidEmail = newAppError(
		CodeUserInvalidEmail,
		"メールアドレスの形式が正しくありません",
	)

	ErrPasswordTooShort = newAppError(
		CodeUserPasswordTooShort,
		"パスワードは8文字以上入力してください",
	)

	ErrInvalidPassword = newAppError(
		CodeUserInvalidPassword,
		"パスワードは英字・数字・記号を含めてください",
	)

	ErrUserNotFound = newAppError(
		CodeUserNotFound,
		"ユーザーが見つかりません",
	)

	ErrEmailAlreadyExists = newAppError(
		CodeUserEmailAlreadyExists,
		"このメールアドレスは既に登録されています",
	)

	ErrNickNameAlreadyTaken = newAppError(
		CodeUserNickNameAlreadyTaken,
		"このニックネームは既に使われています",
	)

	ErrInvalidLogin = newAppError(
		CodeUserInvalidLogin,
		"メール or パスワードが違います",
	)
)
