PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS "profile" (
  "id" INTEGER PRIMARY KEY,
  "display_name" TEXT NOT NULL UNIQUE,
  "avatar" TEXT DEFAULT NULL,
  "description" TEXT DEFAULT NULL,
  "is_public" BOOLEAN DEFAULT true,
  "is_person" BOOLEAN NOT NULL,
  "created_at" DATETIME DEFAULT (CURRENT_TIMESTAMP) NOT NULL
);

CREATE TABLE IF NOT EXISTS "person" (
  "id" integer PRIMARY KEY,
  "email" TEXT UNIQUE NOT NULL,
  "first_name" TEXT NOT NULL,
  "last_name" TEXT NOT NULL,
  "password_hash" TEXT NOT NULL,
  "date_of_birth" TEXT,
  "gender" TEXT NOT NULL,
  FOREIGN KEY ("id") REFERENCES "profile" ("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "group" (
  "id" INTEGER PRIMARY KEY,
  "creator_id" INTEGER NOT NULL,
  FOREIGN KEY ("id") REFERENCES "profile" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("creator_id") REFERENCES "person" ("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "groupmember" (
  "id" INTEGER PRIMARY KEY,
  "group_id" INTEGER NOT NULL,
  "person_id" INTEGER NOT NULL,
  "active" INTEGER DEFAULT 0,
  FOREIGN KEY ("group_id") REFERENCES "group" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("person_id") REFERENCES "person" ("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "follow" (
  "follower_id" INTEGER NOT NULL,
  "following_id" INTEGER NOT NULL,
  "status" INTEGER NOT NULL,
  -- 0: pending, 1: following, 2: blocked
  CHECK (follower_id <> following_id),
  PRIMARY KEY (follower_id, following_id),
  FOREIGN KEY ("follower_id") REFERENCES "person" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("following_id") REFERENCES "profile" ("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "comments" (
  "id" INTEGER PRIMARY KEY,
  "post_id" INTEGER NOT NULL,
  "user_id" INTEGER NOT NULL,
  "content" TEXT NOT NULL,
  "image_path" TEXT,
  "created_at" DATETIME DEFAULT (CURRENT_TIMESTAMP),
  FOREIGN KEY ("post_id") REFERENCES "posts" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("user_id") REFERENCES "person" ("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "message" (
  "sender_id" INTEGER NOT NULL,
  "receiver_id" INTEGER NOT NULL,
  "isread" BOOLEAN DEFAULT 0,
  "content" TEXT NOT NULL,
  -- TODO: in case of notification content should be JSON.
  "created_at" DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  PRIMARY KEY ("sender_id", "receiver_id", "created_at"),
  FOREIGN KEY ("sender_id") REFERENCES "person" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("receiver_id") REFERENCES "profile" ("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "request" (
  "sender_id" INTEGER NOT NULL,
  "receiver_id" INTEGER NOT NULL,
  "is_accept" BOOLEAN DEFAULT 0 CHECK (type IN (0, 1)),
  "type" INTEGER NOT NULL DEFAULT 0 CHECK (type IN (0, 1)),
  "created_at" DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  PRIMARY KEY ("sender_id", "receiver_id", "created_at"),
  FOREIGN KEY ("sender_id") REFERENCES "person" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("receiver_id") REFERENCES "profile" ("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "posts" (
  "id" INTEGER PRIMARY KEY,
  "user_id" INTEGER NOT NULL,
  "group_id" INTEGER,
  "content" TEXT NOT NULL,
  "image_path" TEXT DEFAULT null,
  "privacy" TEXT NOT NULL DEFAULT "public",
  "created_at" DATETIME DEFAULT (CURRENT_TIMESTAMP),
  FOREIGN KEY ("user_id") REFERENCES "person" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("group_id") REFERENCES "group" ("id") ON DELETE CASCADE
);

-- CREATE TABLE IF NOT EXISTS likes (
--   id INTEGER PRIMARY KEY AUTOINCREMENT,
--   user_id INTEGER NOT NULL,
--   post_id INTEGER,
--   comment_id INTEGER,
--   created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
--   FOREIGN KEY (user_id) REFERENCES person(id) ON DELETE CASCADE,
--   FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
--   FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
--   CHECK (
--     (
--       post_id IS NOT NULL
--       AND comment_id IS NULL
--     )
--     OR (
--       post_id IS NULL
--       AND comment_id IS NOT NULL
--     )
--   ),
--   UNIQUE (user_id, post_id, comment_id) -- prevents duplicate likes
-- );