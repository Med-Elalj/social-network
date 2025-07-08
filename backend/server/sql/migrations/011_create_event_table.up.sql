CREATE TABLE IF NOT EXISTS "events" (
  "id" INTEGER PRIMARY KEY,
  "user_id" INTEGER NOT NULL, --event creator
  "group_id" INTEGER, --group event
  "desc" TEXT NOT NULL, --description
  "title" TEXT DEFAULT NULL, --Title
  "timeof" DATETIME NOT NULL, --time of event
  "created_at" DATETIME DEFAULT (CURRENT_TIMESTAMP), --timeofevent
  FOREIGN KEY ("user_id") REFERENCES "user" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("group_id") REFERENCES "group" ("id") ON DELETE CASCADE
);
