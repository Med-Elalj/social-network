CREATE TABLE IF NOT EXISTS "groupmember" (
  "id" INTEGER PRIMARY KEY,
  "group_id" INTEGER NOT NULL,
  "user_id" INTEGER NOT NULL,
  "active" INTEGER DEFAULT 0,
  FOREIGN KEY ("group_id") REFERENCES "group" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON DELETE CASCADE
);
