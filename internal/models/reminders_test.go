package models

import (
	"testing"

	"github.com/tg2648/grem/internal/assert"
)

func TestReminderModelGet(t *testing.T) {
	tests := []struct {
		name       string
		reminderId int
		want       bool
	}{
		{
			name:       "Valid ID",
			reminderId: 1,
			want:       true,
		},
		{
			name:       "Zero ID",
			reminderId: 0,
			want:       false,
		},
		{
			name:       "Non-existent ID",
			reminderId: 2,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)

			m := ReminderModel{db}

			reminder, _ := m.Get(tt.reminderId)
			assert.Equal(t, reminder != nil, tt.want)
		})
	}
}
