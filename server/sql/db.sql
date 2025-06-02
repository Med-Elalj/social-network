PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS "profile" (
  "id" INTEGER PRIMARY KEY,
  "display_name" TEXT Not NULL UNIQUE,
  "avatar" TEXT UNIQUE DEFAULT null,
  "description" TEXT DEFAULT null,
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
  "date_of_birth" TEXT NOT NULL,
  "gender" INTEGER NOT NULL, -- 0: male, 1: female, 2: Attack Helicopter
  FOREIGN KEY ("id") REFERENCES "profile" ("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "group" (
  "id" INTEGER PRIMARY KEY,
  "creator_id" INTEGER NOT NULL,
  FOREIGN KEY ("id") REFERENCES "profile" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("creator_id") REFERENCES "person" ("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "posts" (
  "id" INTEGER PRIMARY KEY,
  "user_id" INTEGER NOT NULL,
  "group_id" INTEGER,
  "title" TEXT NOT NULL,
  "content" TEXT NOT NULL,
  "image_path" TEXT UNIQUE DEFAULT null,
  "privacy" INTEGER NOT NULL DEFAULT 0, -- 0: public, 1: followers/group members, 2: private
  "created_at" DATETIME DEFAULT (CURRENT_TIMESTAMP),
  FOREIGN KEY ("user_id") REFERENCES "person" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("group_id") REFERENCES "group" ("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "follow" (
  "follower_id" INTEGER NOT NULL,
  "following_id" INTEGER NOT NULL,
  "status" INTEGER NOT NULL,  -- 0: pending, 1: following, 2: blocked
  CHECK (follower_id <> following_id),
  PRIMARY KEY (follower_id,following_id),
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
  "content" TEXT NOT NULL, -- TODO: in case of notification content should be JSON.
  "created_at" DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  PRIMARY KEY ("sender_id", "receiver_id", "created_at"),
  FOREIGN KEY ("sender_id") REFERENCES "person" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("receiver_id") REFERENCES "profile" ("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_files (
  uid INTEGER NOT NULL,
  filename TEXT NOT NULL,
  size INTEGER NOT NULL,
  PRIMARY KEY (uid, filename),
  FOREIGN KEY (uid) REFERENCES person (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS pivate_post_visibility (
  post_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  PRIMARY KEY (post_id, user_id),
  FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES person (id) ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS prevent_user_storage_over_10mb
BEFORE INSERT ON user_files
FOR EACH ROW
BEGIN
  SELECT
    CASE
      WHEN (
        (SELECT IFNULL(SUM(size),0) FROM user_files WHERE uid = NEW.uid) + NEW.size
      ) > 10485760
      THEN
        RAISE(ABORT, 'User storage quota exceeded (10MB limit)')
    END;
END;

CREATE TRIGGER IF NOT EXISTS check_event_post_group
BEFORE INSERT ON posts
FOR EACH ROW
WHEN NEW.title = '<EVENT>' AND NEW.group_id IS NULL
BEGIN
  SELECT RAISE(ABORT, 'Event posts must have a group_id');
END;