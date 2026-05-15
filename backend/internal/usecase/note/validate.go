package note

import (
	"strings"
	"unicode/utf8"

	"notes-app/internal/apperror"
)

const (
	maxTitleLength   = 50
	maxContentLength = 1000
)

// タイトルチェック
func validateTitle(title string) error {
	title = strings.TrimSpace(title)

	if title == "" {
		return apperror.ErrTitleRequired
	}

	if utf8.RuneCountInString(title) > maxTitleLength {
		return apperror.ErrTitleTooLong
	}

	return nil
}

// 本文チェック
func validateContent(content string) error {
	content = strings.TrimSpace(content)

	if content == "" {
		return apperror.ErrContentRequired
	}

	if utf8.RuneCountInString(content) > maxContentLength {
		return apperror.ErrContentTooLong
	}

	return nil
}
