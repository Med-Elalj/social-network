CREATE TABLE IF NOT EXISTS "profile" (
  "id" INTEGER PRIMARY KEY,
  "email" TEXT UNIQUE NOT NULL,
  "first_name" TEXT NOT NULL,
  "last_name" TEXT NOT NULL,
  "display_name" TEXT NOT NULL UNIQUE,
  "date_of_birth" TEXT,
  "gender" TEXT NOT NULL,
  "avatar" TEXT DEFAULT NULL,
  "description" TEXT DEFAULT NULL,
  "is_public" BOOLEAN DEFAULT true,
  "is_user" BOOLEAN NOT NULL,
  "created_at" DATETIME DEFAULT (CURRENT_TIMESTAMP) NOT NULL
);