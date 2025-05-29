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
  "gender" INTEGER NOT NULL,
  FOREIGN KEY ("id") REFERENCES "profile" ("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "group" (
  "id" INTEGER PRIMARY KEY,
  "creator_id" INTEGER NOT NULL,
  FOREIGN KEY ("id") REFERENCES "profile" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("creator_id") REFERENCES "profile" ("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "posts" (
  "id" INTEGER PRIMARY KEY,
  "user_id" INTEGER NOT NULL,
  "group_id" INTEGER,
  "title" TEXT NOT NULL,
  "content" TEXT NOT NULL,
  "image_path" TEXT UNIQUE DEFAULT null,
  "privacy" TEXT NOT NULL DEFAULT 'public',
  "created_at" DATETIME DEFAULT (CURRENT_TIMESTAMP),
  FOREIGN KEY ("user_id") REFERENCES "person" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("group_id") REFERENCES "group" ("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "follow" (
  "follower_id" INTEGER NOT NULL,
  "following_id" INTEGER NOT NULL,
  "status" INTEGER NOT NULL,
  CHECK (follower_id <> following_id),
  PRIMARY KEY (follower_id,following_id),
  FOREIGN KEY ("follower_id") REFERENCES "profile" ("id") ON DELETE CASCADE,
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
  FOREIGN KEY ("user_id") REFERENCES "profile" ("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "post_interactions" (
  "user_id" INTEGER,
  "post_id" INTEGER,
  "interaction" BOOLEAN NOT NULL,
  PRIMARY KEY ("user_id", "post_id"),
  FOREIGN KEY ("post_id") REFERENCES "posts" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("user_id") REFERENCES "profile" ("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "group_events" (
  "id" INTEGER PRIMARY KEY,
  "group_id" INTEGER NOT NULL,
  "title" TEXT NOT NULL,
  "description" TEXT NOT NULL,
  "event_date" DATETIME NOT NULL,
  "creator_id" INTEGER NOT NULL,
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  "options" TEXT DEFAULT '["Going","Not Going"]',
  FOREIGN KEY ("creator_id") REFERENCES "profile" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("group_id") REFERENCES "group" ("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "eventResponse" (
  "event_id" INTEGER NOT NULL,
  "user_id" INTEGER NOT NULL,
  "response" INTEGER NOT NULL,
  "created_at" DATETIME DEFAULT (CURRENT_TIMESTAMP),
  PRIMARY KEY ("event_id", "user_id"),
  FOREIGN KEY ("event_id") REFERENCES "group_events" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("user_id") REFERENCES "person" ("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "message" (
  "sender_id" INTEGER NOT NULL,
  "recever_id" INTEGER NOT NULL,
  "isread" BOOLEAN DEFAULT 0,
  "content" TEXT NOT NULL,
  "created_at" DATETIME DEFAULT (CURRENT_TIMESTAMP),
  PRIMARY KEY ("sender_id", "recever_id", "created_at")
  FOREIGN KEY ("sender_id") REFERENCES "profile" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("recever_id") REFERENCES "profile" ("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "notifications" (
  "id" INTEGER PRIMARY KEY,
  "type" TEXT NOT NULL,
  "related_id" INTEGER,
  "recever_id" INTEGER NOT NULL,
  "sender_id" INTEGER NOT NULL,
  "is_read" BOOLEAN DEFAULT 0,
  "created_at" DATETIME DEFAULT (CURRENT_TIMESTAMP),
  FOREIGN KEY ("recever_id") REFERENCES "profile" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("sender_id") REFERENCES "profile" ("id") ON DELETE CASCADE
);

-- add person transaction
-- BEGIN TRANSACTION;

-- Step 1: Insert into profile
-- INSERT INTO profile (display_name, avatar, description, is_public, is_person)
-- VALUES ('jdoe', null, 'Just a sample user', 1, 1);

-- -- Step 2: Insert into person using last inserted profile id
-- INSERT INTO person (ent, email, first_name, last_name, password_hash, date_of_birth)
-- VALUES (last_insert_rowid(), 'jdoe@example.com', 'John', 'Doe', 'hashed_password_here', '1990-01-01');

-- COMMIT;

-- -- add group transaction
-- BEGIN TRANSACTION;

-- INSERT INTO "profile" (display_name, avatar, description, is_public, is_person)
-- VALUES ('GroupName', null, 'Just a sample group', 1, 0);

-- INSERT INTO "group" (id, creator_id)
-- VALUES (last_insert_rowid(), 'id_of_creator_profile');

-- COMMIT;