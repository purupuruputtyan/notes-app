package user

// このファイルは PostgreSQL / lib/pq 連携で使う識別子をまとめる。
const (
	// pgSQLStateUniqueViolation は unique_violation（一意制約違反）。
	pgSQLStateUniqueViolation = "23505"
)

// マイグレーションの CREATE TABLE ... UNIQUE に付く既定の制約名。
const (
	constraintUsersEmailKey    = "users_email_key"
	constraintUsersNickNameKey = "users_nick_name_key"
)
