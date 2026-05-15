package note

import (
	"context"

	"notes-app/internal/models"
)

// StubRepo はユースケース／ハンドラのテスト用インメモリ実装（本番コードからは使わないこと）。
type StubRepo struct {
	notes    models.NoteSlice
	indexErr error // 非 nil のとき Index は常にこのエラーを返す（リポジトリ障害のシミュレーション）
}

func (s *StubRepo) Index(_ context.Context, userID string) (models.NoteSlice, error) {
	if s.indexErr != nil {
		return nil, s.indexErr
	}
	var out models.NoteSlice
	for _, n := range s.notes {
		if n.UserID == userID {
			out = append(out, n)
		}
	}
	return out, nil
}
