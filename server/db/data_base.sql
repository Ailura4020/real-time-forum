CREATE TABLE "USERS" (
	"UserId"	INTEGER PRIMARY KEY AUTOINCREMENT,
	"Nickname"	TEXT NOT NULL UNIQUE,
	"Age"	INTEGER,
	"Gender"	TEXT CHECK ( Gender IN ( 'male', 'female', 'other' )),
	"FirstName"	TEXT NOT NULL,
	"LastName"	TEXT NOT NULL,
	"Email"	TEXT NOT NULL UNIQUE,
	"Password"	TEXT NOT NULL,
	"DateRegister"	TEXT NOT NULL,
--     "Role" TEXT NOT NULL DEFAULT 'user',
	"Status" TEXT DEFAULT 'offline'
);
--
-- CREATE TABLE "POSTS" (
-- 	"PostId"	INTEGER PRIMARY KEY AUTOINCREMENT,
-- 	"Title"	TEXT NOT NULL,
-- 	"Category"	TEXT NOT NULL,
-- 	"TextContent"	TEXT,
-- 	"DateCreation"	DATETIME DEFAULT CURRENT_TIMESTAMP,
-- 	"UserId"	INTEGER,
-- 	FOREIGN KEY ("UserId") REFERENCES "USERS"("UserId")
-- );
--
-- CREATE TABLE "COMMENTS" (
--     "CommentsId" INTEGER PRIMARY KEY AUTOINCREMENT,
--     "TextContent" TEXT NOT NULL,
--     "CreateDate" DATETIME DEFAULT CURRENT_TIMESTAMP,
--     "UserId" INTEGER,
--     "PostId" INTEGER,
--     FOREIGN KEY ("PostId") REFERENCES "POSTS"("PostId"),
--     FOREIGN KEY ("UserId") REFERENCES "USERS"("UserId")
-- );
--
-- CREATE TABLE "PRIVATEMESSAGE" (
--     "PrivateMessageId" INTEGER PRIMARY KEY AUTOINCREMENT,
--     "TextContent" TEXT NOT NULL,
--     "DateSent" DATETIME DEFAULT CURRENT_TIMESTAMP,
--     "SenderId" INTEGER,
--     "ReceiverId" INTEGER,
--     FOREIGN KEY ("SenderId") REFERENCES "USERS"("UserId"),
--     FOREIGN KEY ("ReceiverId") REFERENCES "USERS"("UserId")
-- );
--
-- CREATE TABLE "NOTIF" (
--
-- )

-- NOTIF
-- LIKE/DISLIKE