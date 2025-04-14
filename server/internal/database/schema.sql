CREATE TABLE IF NOT EXISTS users (
    id   INTEGER PRIMARY KEY,
    userName TEXT NOT NULL UNIQUE,
    password  TEXT NOT NULL,
    created_at INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS sessions (
    token char(43) primary key,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);
CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions (expiry);
