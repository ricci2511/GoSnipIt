package models

import (
	"testing"

	"gosnipit.ricci2511.dev/internal/assert"
)

func TestUserModelExists(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Valid userID",
			userID: 1,
			want:   true,
		},
		{
			name:   "Zero userID",
			userID: 0,
			want:   false,
		},
		{
			name:   "Non-existent userID",
			userID: 2,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// each test case sets up a clean instance of the test database
			db := newTestDb(t)

			m := UserModel{db}

			exists, err := m.Exists(tt.userID)

			// check if the user exists, in case of error test that the error is nil
			assert.Equal(t, exists, tt.want)
			assert.NilError(t, err)
		})
	}
}
