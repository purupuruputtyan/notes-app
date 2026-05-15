package httputil

import (
	"errors"
	"net/http"
	"testing"

	"notes-app/internal/apperror"
)

func TestClientStatusFromAppError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "user not found",
			err:        apperror.ErrUserNotFound,
			wantStatus: http.StatusNotFound,
			wantMsg:    "not found",
		},
		{
			name:       "note owner not found",
			err:        apperror.ErrOwnerNotFound,
			wantStatus: http.StatusNotFound,
			wantMsg:    "not found",
		},
		{
			name:       "validation",
			err:        apperror.ErrInvalidEmail,
			wantStatus: http.StatusBadRequest,
			wantMsg:    apperror.ErrInvalidEmail.Message,
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
			st, msg := ClientStatusFromAppError(tt.err)
			if st != tt.wantStatus || msg != tt.wantMsg {
				t.Fatalf("got (%d, %q), want (%d, %q)", st, msg, tt.wantStatus, tt.wantMsg)
			}
		})
	}
}
