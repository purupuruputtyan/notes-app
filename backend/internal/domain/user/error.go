package user

type ValidateError struct {
	Code    string
	Message string
}

func (e ValidateError) Error() string {
	return e.Message
}

var (
	ErrNickNameRequired = ValidateError{
		Code:    "ERR_001",
		Message: "ニックネームを入力してください",
	}

	ErrNickNameTooLong = ValidateError{
		Code:    "ERR_002",
		Message: "ニックネームは20文字以内で入力してください",
	}

	ErrInvalidEmail = ValidateError{
		Code:    "ERR_003",
		Message: "メールアドレスの形式が正しくありません",
	}

	ErrPasswordTooShort = ValidateError{
		Code:    "ERR_004",
		Message: "パスワードは8文字以上入力してください",
	}

	ErrInvalidPassword = ValidateError{
		Code:    "ERR_005",
		Message: "パスワードは英字・数字・記号を含めてください",
	}
)
