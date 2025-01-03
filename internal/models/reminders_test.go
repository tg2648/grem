package models

import (
	"testing"
	"time"

	"github.com/tg2648/grem/internal/assert"
)

func TestReminderModelGet(t *testing.T) {
	tests := []struct {
		name        string
		reminderId  int
		want        bool
		expectedErr error
	}{
		{
			name:        "Valid ID",
			reminderId:  1,
			want:        true,
			expectedErr: nil,
		},
		{
			name:        "Zero ID",
			reminderId:  0,
			want:        false,
			expectedErr: ErrNoRecord,
		},
		{
			name:        "Non-existent ID",
			reminderId:  999,
			want:        false,
			expectedErr: ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)

			m := ReminderModel{db}

			reminder, err := m.Get(tt.reminderId)
			assert.Equal(t, reminder != nil, tt.want)
			assert.Equal(t, err, tt.expectedErr)
		})
	}
}

func TestReminderModelGetDue(t *testing.T) {
	t.Run("One result", func(t *testing.T) {
		db := newTestDB(t)
		m := ReminderModel{db}

		due := time.Date(2024, time.December, 28, 0, 0, 0, 0, time.UTC)
		reminders, err := m.GetDue(&due)

		assert.Equal(t, len(reminders), 1)
		assert.NilError(t, err)
	})

	t.Run("Zero results", func(t *testing.T) {
		db := newTestDB(t)
		m := ReminderModel{db}

		due := time.Date(2023, time.December, 28, 0, 0, 0, 0, time.UTC)
		reminders, err := m.GetDue(&due)

		assert.Equal(t, len(reminders), 0)
		assert.NilError(t, err)
	})
}
