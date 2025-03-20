-- CREATE INDEX idx_comments_userid ON COMMENTS(UserId);
-- CREATE INDEX idx_comments_postid ON COMMENTS(PostId);
--
-- CREATE TABLE "USERS"
-- (
--     "UserId"       INTEGER PRIMARY KEY AUTOINCREMENT,
--     "Nickname"     TEXT NOT NULL UNIQUE,
--     "Age"          INTEGER,
--     "Gender"       TEXT CHECK ( Gender IN ('male', 'female', 'other')),
--     "FirstName"    TEXT NOT NULL,
--     "LastName"     TEXT NOT NULL,
--     "Email"        TEXT NOT NULL UNIQUE,
--     "Password"     TEXT NOT NULL,
--     "DateRegister" TEXT NOT NULL,
-- --     "Role" TEXT NOT NULL DEFAULT 'user',
--     "Status"       TEXT DEFAULT 'offline'
-- );
-- --
-- CREATE TABLE "POSTS"
-- (
--     "PostId"       INTEGER PRIMARY KEY AUTOINCREMENT,
--     "Title"        TEXT NOT NULL,
--     "Category"     TEXT NOT NULL,
--     "TextContent"  TEXT,
--     "DateCreation" DATETIME DEFAULT CURRENT_TIMESTAMP,
--     "UserId"       INTEGER,
--     "Likes"        INTEGER  DEFAULT 0,
--     "Dislikes"     INTEGER  DEFAULT 0,
--     FOREIGN KEY ("UserId") REFERENCES "USERS" ("UserId")
-- );
--
-- CREATE TABLE "COMMENTS"
-- (
--     "CommentsId"  INTEGER PRIMARY KEY AUTOINCREMENT,
--     "TextContent" TEXT NOT NULL,
--     "CreateDate"  DATETIME DEFAULT CURRENT_TIMESTAMP,
--     "UserId"      INTEGER,
--     "PostId"      INTEGER,
--     "Likes"       INTEGER  DEFAULT 0,
--     "Dislikes"    INTEGER  DEFAULT 0,
--     FOREIGN KEY ("PostId") REFERENCES "POSTS" ("PostId"),
--     FOREIGN KEY ("UserId") REFERENCES "USERS" ("UserId")
-- );

CREATE TABLE "PRIVATEMESSAGE"
(
    "PrivateMessageId" INTEGER PRIMARY KEY AUTOINCREMENT,
    "TextContent"      TEXT NOT NULL,
    "DateSent"         DATETIME DEFAULT CURRENT_TIMESTAMP,
    "SenderId"         INTEGER,
    "ReceiverId"       INTEGER,
    FOREIGN KEY ("SenderId") REFERENCES "USERS" ("UserId"),
    FOREIGN KEY ("ReceiverId") REFERENCES "USERS" ("UserId")
);

CREATE TABLE "NOTIF"
(
    "NotifId"     INTEGER PRIMARY KEY AUTOINCREMENT,
    "UserId"      INTEGER NOT NULL,
    "Message"     TEXT    NOT NULL,
    "DateCreated" DATETIME DEFAULT CURRENT_TIMESTAMP,
    "IsRead"      BOOLEAN  DEFAULT 0,
    FOREIGN KEY ("UserId") REFERENCES "USERS" ("UserId")
)

-- INSERT INTO NOTIF (UserId, Message) VALUE (1, 'Vous avez un nouveau message.');
-- INSERT INTO NOTIF (UserId, Message) VALUE (2, 'Votre post a été aimé.');
-- NOTIF
-- LIKE/DISLIKE