CREATE TABLE "users"
(
    "user_id"       INTEGER PRIMARY KEY AUTOINCREMENT,
    "nickname"      TEXT NOT NULL UNIQUE,
    "age"           INTEGER CHECK (age >= 0),
    "gender"        TEXT CHECK (gender IN ('male', 'female', 'other')),
    "first_name"    TEXT NOT NULL,
    "last_name"     TEXT NOT NULL,
    "email"         TEXT NOT NULL UNIQUE CHECK (email LIKE '%@%.%'),
    "password"      TEXT NOT NULL,
    "date_register" DATETIME DEFAULT CURRENT_TIMESTAMP,
    "status"        BOOLEAN  DEFAULT FALSE
);

CREATE TABLE "posts"
(
    "post_id"       INTEGER PRIMARY KEY AUTOINCREMENT,
    "title"         TEXT    NOT NULL,
    "category"      TEXT    NOT NULL,
    "text_content"  TEXT    NOT NULL,
    "date_creation" DATETIME DEFAULT CURRENT_TIMESTAMP,
    "user_id"       INTEGER NOT NULL,
    "likes"         INTEGER  DEFAULT 0 CHECK (likes >= 0),
    "dislikes"      INTEGER  DEFAULT 0 CHECK (dislikes >= 0),
    FOREIGN KEY ("user_id") REFERENCES "users" ("user_id") ON DELETE CASCADE
);

CREATE TABLE "comments"
(
    "comments_id"  INTEGER PRIMARY KEY AUTOINCREMENT,
    "text_content" TEXT    NOT NULL,
    "create_date"  DATETIME DEFAULT CURRENT_TIMESTAMP,
    "user_id"      INTEGER NOT NULL,
    "post_id"      INTEGER NOT NULL,
    "likes"        INTEGER  DEFAULT 0 CHECK (likes >= 0),
    "dislikes"     INTEGER  DEFAULT 0 CHECK (dislikes >= 0),
    FOREIGN KEY ("post_id") REFERENCES "posts" ("post_id") ON DELETE CASCADE,
    FOREIGN KEY ("user_id") REFERENCES "users" ("user_id") ON DELETE CASCADE
);

CREATE INDEX idx_comments_user_id ON comments (user_id);
CREATE INDEX idx_comments_post_id ON comments (post_id);
CREATE INDEX idx_posts_user_id ON posts (user_id);
CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_nickname ON users (nickname);

CREATE TABLE IF NOT EXISTS "private_messages"
(
    "id"
        INTEGER
        PRIMARY
            KEY
        AUTOINCREMENT,
    "sender_id"
        INTEGER
        NOT
            NULL,
    "receiver_id"
        INTEGER
        NOT
            NULL,
    "content"
        TEXT
        NOT
            NULL,
    "timestamp"
        DATETIME
        DEFAULT
            CURRENT_TIMESTAMP,
    FOREIGN
        KEY
        (
         "sender_id"
            ) REFERENCES "users"
        (
         "user_id"
            ) ON DELETE CASCADE,
    FOREIGN KEY
        (
         "receiver_id"
            ) REFERENCES "users"
        (
         "user_id"
            )
        ON DELETE CASCADE
);

CREATE INDEX idx_messages_sender ON private_messages (sender_id);
CREATE INDEX idx_messages_receiver ON private_messages (receiver_id);