package models

import (
	"testing"

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
			reminderId:  2,
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
	
}