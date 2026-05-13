package user

type AppError struct {
	Code    string
	Message string
}

func (e AppError) Error() string {
	return e.Message
}

var (
	ErrNickNameRequired = AppError{
		Code:    "ERR_001",
		Message: "ニックネームを入力してください",
	}

	ErrNickNameTooLong = AppError{
		Code:    "ERR_002",
		Message: "ニックネームは20文字以内で入力してください",
	}

	ErrInvalidEmail = AppError{
		Code:    "ERR_003",
		Message: "メールアドレスの形式が正しくありません",
	}

	ErrPasswordTooShort = AppError{
		Code:    "ERR_004",
		Message: "パスワードは8文字以上入力してください",
	}

	ErrInvalidPassword = AppError{
		Code:    "ERR_005",
		Message: "パスワードは英字・数字・記号を含めてください",
	}

	ErrUserNotFound = AppError{
		Code:    "ERR_006",
		Message: "ユーザーが見つかりません",
	}
)
