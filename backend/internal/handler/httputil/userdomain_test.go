package httputil

import (
	"errors"
	"net/http"
	"testing"

	domain "notes-app/internal/domain/user"
)

func TestClientStatusFromUserDomain(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "not found",
			err:        domain.ErrUserNotFound,
			wantStatus: http.StatusNotFound,
			wantMsg:    "not found",
		},
		{
			name:       "validation",
			err:        domain.ErrInvalidEmail,
			wantStatus: http.StatusBadRequest,
			wantMsg:    domain.ErrInvalidEmail.Message,
		},
		{
			name:       "internal",
			err:        errors.New("db exploded"),
			wantStatus: http.StatusInternalServerError,
			wantMsg:    "internal server error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			st, msg := ClientStatusFromUserDomain(tt.err)
			if st != tt.wantStatus || msg != tt.wantMsg {
				t.Fatalf("got (%d, %q), want (%d, %q)", st, msg, tt.wantStatus, tt.wantMsg)
			}
		})
	}
}
