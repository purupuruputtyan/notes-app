package user

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/aarondl/null/v8"

	"notes-app/internal/models"
)

func TestUserRepository_Create(t *testing.T) {
	db := openTestDB(t)
	repo := NewUser(db)

	user := models.User{
		NickName: fmt.Sprintf(
			"テストユーザー-%d",
			time.Now().UnixNano(),
		),
		Email: fmt.Sprintf(
			"test-%d@example.com",
			time.Now().UnixNano(),
		),
		PasswordHash: "52fdb484dce77487c673a0854efe360a212ac44b3deeda0eb9b33777e0a6e11e",
		IconImage:    null.StringFrom("https://example.com"),
	}

	created, err := repo.Create(user)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// テストで作ったデータを削除
	t.Cleanup(func() {
		_, err := db.Exec(
			"DELETE FROM users WHERE id = $1",
			created.ID,
		)

		if err != nil {
			t.Fatalf(
				"failed to cleanup user: %v",
				err,
			)
		}
	})

	if created.ID == "" {
		t.Fatalf("expected id to be set")
	}

	if created.NickName != user.NickName {
		t.Fatalf(
			"expected NickName %s, got %s",
			user.NickName,
			created.NickName,
		)
	}

	if created.Email != user.Email {
		t.Fatalf(
			"expected Email %s, got %s",
			user.Email,
			created.Email,
		)
	}

	if created.PasswordHash != user.PasswordHash {
		t.Fatalf(
			"expected PasswordHash %s, got %s",
			user.PasswordHash,
			created.PasswordHash,
		)
	}

	if !created.IsActive {
		t.Fatalf("expected IsActive to be true")
	}
}

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		envOrDefault("DB_HOST", "localhost"),
		envOrDefault("DB_PORT", "5432"),
		envOrDefault("DB_USER", "notes"),
		envOrDefault("DB_PASSWORD", "notes"),
		envOrDefault("DB_NAME", "notes_app"),
		envOrDefault("DB_SSLMODE", "disable"),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("failed to connect db: %v", err)
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	return db
}

func envOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}
