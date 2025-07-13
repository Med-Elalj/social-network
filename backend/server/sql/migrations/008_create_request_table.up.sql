CREATE TABLE IF NOT EXISTS "request" (
  "sender_id" INTEGER NOT NULL,
  "receiver_id" INTEGER NOT NULL,
  "target_id" INTEGER, -- optional: group.id or event.id depending on type
  "type" INTEGER NOT NULL DEFAULT 0 CHECK (type IN (0, 1, 2)), -- 0=user, 1=group, 2=event
  "created_at" DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  PRIMARY KEY ("sender_id", "receiver_id", "created_at"),
  FOREIGN KEY ("sender_id") REFERENCES "profile" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("receiver_id") REFERENCES "profile" ("id") ON DELETE CASCADE
);
