CREATE TABLE IF NOT EXISTS "request" (
  "sender_id" INTEGER NOT NULL,
  "receiver_id" INTEGER NOT NULL,
  "is_accept" BOOLEAN DEFAULT 0 CHECK (type IN (0, 1)),
  "type" INTEGER NOT NULL DEFAULT 0 CHECK (type IN (0, 1)),
  "created_at" DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  PRIMARY KEY ("sender_id", "receiver_id", "created_at"),
  FOREIGN KEY ("sender_id") REFERENCES "user" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("receiver_id") REFERENCES "profile" ("id") ON DELETE CASCADE
);
