CREATE TABLE "users" (
  "id" INTEGER PRIMARY KEY,
  "email" TEXT UNIQUE NOT NULL,
  "nickname" TEXT DEFAULT null,
  "first_name" TEXT NOT NULL,
  "last_name" TEXT NOT NULL,
  "password_hash" TEXT NOT NULL,
  "date_of_birth" TEXT NOT NULL,
  "avatar" TEXT UNIQUE DEFAULT null,
  "about_me" TEXT,
  "is_public" BOOLEAN DEFAULT true,
  "created_at" DATETIME DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "posts" (
  "id" INTEGER PRIMARY KEY,
  "user_id" INTEGER NOT NULL,
  "group_id" INTEGER,
  "content" TEXT NOT NULL,
  "image_path" TEXT UNIQUE DEFAULT null,
  "privacy" TEXT NOT NULL DEFAULT 'public',
  "created_at" DATETIME DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "followers" (
  "id" INTEGER PRIMARY KEY,
  "follower_id" INTEGER NOT NULL,
  "following_id" INTEGER NOT NULL,
  "status" TEXT NOT NULL DEFAULT 'pending',
  "created_at" DATETIME DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "comments" (
  "id" INTEGER PRIMARY KEY,
  "post_id" INTEGER NOT NULL,
  "user_id" INTEGER NOT NULL,
  "content" TEXT NOT NULL,
  "image_path" TEXT,
  "created_at" DATETIME DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "groups" (
  "id" INTEGER PRIMARY KEY,
  "creator_id" INTEGER NOT NULL,
  "title" TEXT NOT NULL,
  "description" TEXT NOT NULL,
  "created_at" DATETIME DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "post_interactions" (
  "user_id" INTEGER,
  "post_id" INTEGER,
  "interaction" INTEGER DEFAULT 0,
  PRIMARY KEY ("user_id", "post_id")
);

CREATE TABLE "group_events" (
  "id" INTEGER PRIMARY KEY,
  "group_id" INTEGER NOT NULL,
  "title" TEXT NOT NULL,
  "description" TEXT NOT NULL,
  "event_date" DATETIME NOT NULL,
  "creator_id" INTEGER NOT NULL,
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  "options" TEXT DEFAULT '["Going","Not Going"]'
);

CREATE TABLE "eventResponse" (
  "event_id" INTEGER NOT NULL,
  "user_id" INTEGER NOT NULL,
  "response" INTEGER NOT NULL,
  "created_at" DATETIME DEFAULT (CURRENT_TIMESTAMP),
  PRIMARY KEY ("event_id", "user_id")
);

CREATE TABLE "messages" (
  "id" INTEGER PRIMARY KEY,
  "sender_id" INTEGER NOT NULL,
  "recever_id" INTEGER NOT NULL,
  "isread" BOOLEAN NOT NULL DEFAULT 0,
  "content" TEXT NOT NULL,
  "created_at" DATETIME DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE "notifications" (
  "id" INTEGER PRIMARY KEY,
  "type" TEXT NOT NULL,
  "related_id" INTEGER,
  "recever_id" INTEGER NOT NULL,
  "sender_id" INTEGER NOT NULL,
  "is_read" BOOLEAN DEFAULT 0,
  "created_at" DATETIME DEFAULT (CURRENT_TIMESTAMP)
);

ALTER TABLE "posts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "posts" ADD FOREIGN KEY ("group_id") REFERENCES "groups" ("id") ON DELETE CASCADE;

ALTER TABLE "followers" ADD FOREIGN KEY ("follower_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "followers" ADD FOREIGN KEY ("following_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "comments" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id") ON DELETE CASCADE;

ALTER TABLE "comments" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "groups" ADD FOREIGN KEY ("creator_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "post_interactions" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id") ON DELETE CASCADE;

ALTER TABLE "post_interactions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "group_events" ADD FOREIGN KEY ("group_id") REFERENCES "groups" ("id") ON DELETE CASCADE;

ALTER TABLE "group_events" ADD FOREIGN KEY ("creator_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "eventResponse" ADD FOREIGN KEY ("event_id") REFERENCES "group_events" ("id") ON DELETE CASCADE;

ALTER TABLE "eventResponse" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "messages" ADD FOREIGN KEY ("sender_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "messages" ADD FOREIGN KEY ("recever_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "notifications" ADD FOREIGN KEY ("recever_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "notifications" ADD FOREIGN KEY ("sender_id") REFERENCES "users" ("id") ON DELETE CASCADE;
