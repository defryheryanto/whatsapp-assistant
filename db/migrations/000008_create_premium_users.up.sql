CREATE TABLE IF NOT EXISTS premium_users (
    id SERIAL PRIMARY KEY,
    user_jid TEXT NOT NULL
);

CREATE INDEX idx_premium_users_user_jid
ON premium_users (user_jid);
