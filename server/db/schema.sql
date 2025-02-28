-- called schema because it defines the structure of the database (the blueprint of the database.)
CREATE TABLE IF NOT EXISTS USERS
(
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    nickname      TEXT NOT NULL,
    age           INTEGER,
    gender        TEXT,
    first_name    TEXT,
    last_name     TEXT,
    email         TEXT NOT NULL UNIQUE,
    password      TEXT NOT NULL,
    date_register TEXT NOT NULL
);
