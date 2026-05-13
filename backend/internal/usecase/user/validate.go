package user

import (
	"net/mail"
	"regexp"
	"unicode/utf8"

	"notes-app/internal/domain/user"
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
		return user.ErrNickNameRequired
	}
	if utf8.RuneCountInString(nickname) > maxNickNameLength {
		return user.ErrNickNameTooLong
	}
	return nil
}

// メールチェック
func validateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return user.ErrInvalidEmail
	}
	return nil
}

// パスワードチェック（英数字＋記号）
func validatePassword(password string) error {
	if utf8.RuneCountInString(password) < minPasswordLength {
		return user.ErrPasswordTooShort
	}

	hasLetter := passwordHasLetter.MatchString(password)
	hasNumber := passwordHasNumber.MatchString(password)
	hasSymbol := passwordHasSymbol.MatchString(password)

	if !hasLetter || !hasNumber || !hasSymbol {
		return user.ErrInvalidPassword
	}

	return nil
}
