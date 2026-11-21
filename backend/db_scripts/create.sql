CREATE TABLE users (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name TEXT,
    pass TEXT
);

CREATE TABLE tokens (
    refresh_token INTEGER PRIMARY KEY,
    expires_at INTEGER,
    user_id INTEGER,

    FOREIGN KEY (user_id) REFERENCES users (id)
)