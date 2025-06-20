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

CREATE TABLE IF NOT EXISTS "categories" (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  category_name VARCHAR(50) NOT NULL UNIQUE,
  Category_icon_path TEXT NOT NULL UNIQUE
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