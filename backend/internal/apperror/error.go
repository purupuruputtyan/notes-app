package apperror

// AppError はクライアント向け Code / Message を持つアプリケーションエラー。
type AppError struct {
	Code    string
	Message string
}

func (e AppError) Error() string {
	return e.Message
}

func newAppError(code, message string) AppError {
	return AppError{Code: code, Message: message}
}
