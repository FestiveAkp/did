CREATE TABLE activities (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    task_id    INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    note       TEXT    NOT NULL,
    created_at DATETIME NOT NULL DEFAULT (datetime('now'))
);
