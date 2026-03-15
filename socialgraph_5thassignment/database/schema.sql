CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    gender TEXT NOT NULL,
    birth_date DATE NOT NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS user_friends (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    friend_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, friend_id),
    CHECK (user_id <> friend_id)
);