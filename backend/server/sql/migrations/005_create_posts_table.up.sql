CREATE TABLE IF NOT EXISTS "posts" (
  "id" INTEGER PRIMARY KEY,
  "user_id" INTEGER NOT NULL, --event creator
  "group_id" INTEGER, --group event
  "content" TEXT NOT NULL, --description
  "image_path" TEXT DEFAULT NULL, --Title
  "privacy" TEXT NOT NULL DEFAULT "public",
  "created_at" DATETIME DEFAULT (CURRENT_TIMESTAMP), --timeofevent
  FOREIGN KEY ("user_id") REFERENCES "profile" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("group_id") REFERENCES "group" ("id") ON DELETE CASCADE
);
