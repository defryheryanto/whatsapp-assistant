CREATE TABLE IF NOT EXISTS saved_texts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_jid TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL
);

CREATE INDEX idx_saved_texts_group_jid_title
ON saved_texts (group_jid, title);
