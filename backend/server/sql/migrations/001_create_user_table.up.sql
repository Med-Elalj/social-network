CREATE TABLE IF NOT EXISTS "user" (
  "id" INTEGER PRIMARY KEY,
  "email" TEXT UNIQUE NOT NULL,
  "first_name" TEXT NOT NULL,
  "last_name" TEXT NOT NULL,
  "password_hash" TEXT NOT NULL,
  "date_of_birth" TEXT,
  "gender" TEXT NOT NULL,
  FOREIGN KEY ("id") REFERENCES "profile" ("id") ON DELETE CASCADE
);
