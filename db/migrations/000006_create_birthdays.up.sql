CREATE TABLE IF NOT EXISTS birthdays (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    birth_date INTEGER NOT NULL,
    birth_month INTEGER NOT NULL,
    birth_year INTEGER NOT NULL,
    target_chat_jid TEXT NOT NULL
);

CREATE INDEX idx_birthdays_date_month
ON birthdays (birth_date, birth_month);
