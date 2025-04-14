CREATE TABLE "users"
(
    "user_id"       INTEGER PRIMARY KEY AUTOINCREMENT,
    "nickname"      TEXT NOT NULL UNIQUE,
    "age"           INTEGER,
    "gender"        TEXT CHECK (gender IN ('male', 'female', 'other')),
    "first_name"    TEXT NOT NULL,
    "last_name"     TEXT NOT NULL,
    "email"         TEXT NOT NULL UNIQUE,
    "password"      TEXT NOT NULL,
    "date_register" DATETIME DEFAULT CURRENT_TIMESTAMP,  -- Changed to DATETIME for consistency
    -- "status"        TEXT DEFAULT 'offline'
        "status"    BOOLEAN

);

CREATE TABLE "posts"
(
    "post_id"       INTEGER PRIMARY KEY AUTOINCREMENT,
    "title"         TEXT NOT NULL,
    "category"      TEXT NOT NULL,
    "text_content"  TEXT,
    "date_creation" DATETIME DEFAULT CURRENT_TIMESTAMP,
    "user_id"       INTEGER,
    "likes"         INTEGER DEFAULT 0,
    "dislikes"      INTEGER DEFAULT 0,
    FOREIGN KEY ("user_id") REFERENCES "users" ("user_id")
);

CREATE TABLE "comments"
(
    "comments_id"   INTEGER PRIMARY KEY AUTOINCREMENT,
    "text_content"  TEXT NOT NULL,
    "create_date"   DATETIME DEFAULT CURRENT_TIMESTAMP,
    "user_id"       INTEGER,
    "post_id"       INTEGER,
    "likes"         INTEGER DEFAULT 0,
    "dislikes"      INTEGER DEFAULT 0,
    FOREIGN KEY ("post_id") REFERENCES "posts" ("post_id"),
    FOREIGN KEY ("user_id") REFERENCES "users" ("user_id")
);

CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_comments_post_id ON comments(post_id);

CREATE TABLE IF NOT EXISTS private_messages 
(
"id" INTEGER PRIMARY KEY AUTOINCREMENT
"sender_id" INTEGER NOT NULL
"receiver_id" INTEGER NOT NULL
timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
FOREIGN KEY (serder_id) REFERENCES users(user_id)
FOREIGN KEY (receiver_id) REFERENCES users(user_id)
);