CREATE TABLE IF NOT EXISTS users (
    id   INTEGER PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password  TEXT NOT NULL,
    created_at INTEGER DEFAULT (unixepoch('now','localtime')) NOT NULL,
    metal_extractor_lvl INTEGER DEFAULT 1 NOT NULL,
    gas_extractor_lvl INTEGER DEFAULT 1 NOT NULL,
    total_gas_expenses INTEGER DEFAULT 0 NOT NULL,
    total_metal_expenses INTEGER DEFAULT 0 NOT NULL
);

CREATE TABLE IF NOT EXISTS gas_income_history (
    id INTEGER PRIMARY KEY,
    income INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    change_timestamp INTEGER DEFAULT (unixepoch('now','localtime')) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS metal_income_history (
    id INTEGER PRIMARY KEY,
    income INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    change_timestamp INTEGER DEFAULT (unixepoch('now','localtime')) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS sessions (
    token char(43) primary key,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);
CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions (expiry);
