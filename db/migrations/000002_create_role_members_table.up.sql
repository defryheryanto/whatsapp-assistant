CREATE TABLE IF NOT EXISTS role_members (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    role_id INTEGER NOT NULL,
    jid TEXT NOT NULL
);
