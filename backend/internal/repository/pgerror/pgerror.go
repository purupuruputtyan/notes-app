package pgerror

import (
	"errors"

	"github.com/lib/pq"
)

// PostgreSQL SQLSTATE（lib/pq の Code と比較する）。
const (
	SQLStateUniqueViolation     = "23505"
	SQLStateForeignKeyViolation = "23503"
)

// マイグレーションの制約名（PostgreSQL の既定名）。
const (
	ConstraintNotesPkey        = "notes_pkey"
	ConstraintNotesUserIDFkey  = "notes_user_id_fkey"
	ConstraintUsersEmailKey    = "users_email_key"
	ConstraintUsersNickNameKey = "users_nick_name_key"
)

// As は err を *pq.Error に変換できるときその値を返す。
func As(err error) (*pq.Error, bool) {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr, true
	}
	return nil, false
}
