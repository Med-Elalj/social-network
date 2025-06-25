CREATE TABLE IF NOT EXISTS "follow" (
  "follower_id" INTEGER NOT NULL,
  "following_id" INTEGER NOT NULL,
  "status" INTEGER NOT NULL,
  CHECK (follower_id <> following_id),
  PRIMARY KEY (follower_id, following_id),
  FOREIGN KEY ("follower_id") REFERENCES "user" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("following_id") REFERENCES "profile" ("id") ON DELETE CASCADE
);
