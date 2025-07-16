CREATE TABLE IF NOT EXISTS "postrack" (
  "post_id" INTEGER NOT NULL,
  "follower_id" INTEGER NOT NULL,
  PRIMARY KEY ("post_id", "follower_id"),
  FOREIGN KEY ("post_id") REFERENCES "posts" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("follower_id") REFERENCES "profile" ("id") ON DELETE CASCADE
);