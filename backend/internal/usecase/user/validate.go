package user

import (
	"net/mail"
	"regexp"
	"unicode/utf8"

	"notes-app/internal/apperror"
)

const (
	maxNickNameLength = 20
	minPasswordLength = 8
)

var (
	passwordHasLetter = regexp.MustCompile(`[a-zA-Z]`)
	passwordHasNumber = regexp.MustCompile(`[0-9]`)
	passwordHasSymbol = regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]\{\};:'",.<>/?\\|` + "`" + `~]`)
)

// ニックネーム
func validateNickName(nickname string) error {
	if nickname == "" {
		return apperror.ErrNickNameRequired
	}
	if utf8.RuneCountInString(nickname) > maxNickNameLength {
		return apperror.ErrNickNameTooLong
	}
	return nil
}

// メールチェック
func validateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return apperror.ErrInvalidEmail
	}
	return nil
}

// パスワードチェック（英数字＋記号）
func validatePassword(password string) error {
	if utf8.RuneCountInString(password) < minPasswordLength {
		return apperror.ErrPasswordTooShort
	}

	hasLetter := passwordHasLetter.MatchString(password)
	hasNumber := passwordHasNumber.MatchString(password)
	hasSymbol := passwordHasSymbol.MatchString(password)

	if !hasLetter || !hasNumber || !hasSymbol {
		return apperror.ErrInvalidPassword
	}

	return nil
}
