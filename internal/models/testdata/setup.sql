CREATE TABLE IF NOT EXISTS reminders (
    id INTEGER NOT NULL PRIMARY KEY,
    title TEXT NOT NULL,
    due_at DATE NOT NULL,
    dismissed_at DATETIME DEFAULT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO reminders (title, due_at, dismissed_at, created_at) VALUES
("One", "2024-12-30", NULL, "2024-12-28");
