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
