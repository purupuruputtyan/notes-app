package pgerror

import (
	"errors"
	"testing"

	"github.com/lib/pq"
)

func TestAs(t *testing.T) {
	t.Parallel()

	plain := errors.New("plain")

	tests := []struct {
		name   string
		err    error
		wantOK bool
	}{
		{name: "pq error", err: &pq.Error{Code: SQLStateUniqueViolation}, wantOK: true},
		{name: "plain error", err: plain, wantOK: false},
		{name: "nil", err: nil, wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, ok := As(tt.err)
			if ok != tt.wantOK {
				t.Fatalf("ok: want %v, got %v", tt.wantOK, ok)
			}
			if tt.wantOK && got == nil {
				t.Fatal("expected non-nil pq.Error")
			}
		})
	}
}
