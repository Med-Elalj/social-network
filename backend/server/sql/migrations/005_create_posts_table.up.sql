CREATE TABLE IF NOT EXISTS "posts" (
  "id" INTEGER PRIMARY KEY,
  "user_id" INTEGER NOT NULL,
  "group_id" INTEGER,
  "content" TEXT NOT NULL,
  "image_path" TEXT DEFAULT NULL,
  "privacy" TEXT NOT NULL DEFAULT "public",
  "created_at" DATETIME DEFAULT (CURRENT_TIMESTAMP),
  FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("group_id") REFERENCES "group" ("id") ON DELETE CASCADE
);
