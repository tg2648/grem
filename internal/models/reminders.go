package models

import (
	"database/sql"
	"errors"
	"time"
)

type Reminder struct {
	ID          int
	Title       string
	DueAt       *time.Time
	DismissedAt *sql.NullTime
	CreatedAt   *time.Time
}

type ReminderModel struct {
	DB *sql.DB
}

// Insert new reminder into the DB
// Return the ID of the new reminder
func (m *ReminderModel) Insert(title string, dueAt *time.Time) (int, error) {
	stmt := `
		INSERT INTO reminders (title, due_at)
		VALUES(?, ?)
	`
	result, err := m.DB.Exec(stmt, title, dueAt)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Return a specific Reminder by ID
func (m *ReminderModel) Get(id int) (*Reminder, error) {
	stmt := `
		SELECT id, title, due_at, dismissed_at, created_at
		FROM reminders
		WHERE id = ?
	`

	r := &Reminder{}
	err := m.DB.
		QueryRow(stmt, id).
		Scan(&r.ID, &r.Title, &r.DueAt, &r.DismissedAt, &r.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return r, nil
}

// Return reminders due on or before a certain date
func (m *ReminderModel) GetDue(date *time.Time) ([]*Reminder, error) {
	stmt := `
		SELECT id, title, due_at, dismissed_at, created_at
		FROM reminders
		WHERE
			dismissed_at IS NULL AND
			date(due_at) <= date(?)
		ORDER BY
			due_at,
			id
	`

	rows, err := m.DB.Query(stmt, date)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	reminders := []*Reminder{}

	for rows.Next() {
		r := &Reminder{}

		err = rows.Scan(&r.ID, &r.Title, &r.DueAt, &r.DismissedAt, &r.CreatedAt)
		if err != nil {
			return nil, err
		}

		reminders = append(reminders, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reminders, nil
}

// Return reminders due today
func (m *ReminderModel) GetDueToday() ([]*Reminder, error) {
	now := time.Now().UTC()
	return m.GetDue(&now)
}
