CREATE TABLE IF NOT EXISTS "profile" (
  "id" INTEGER PRIMARY KEY,
  "email" TEXT UNIQUE,
  "first_name" TEXT,
  "last_name" TEXT,
  "display_name" TEXT NOT NULL UNIQUE,
  "date_of_birth" TEXT,
  "gender" TEXT,
  "avatar" TEXT DEFAULT NULL,
  "description" TEXT DEFAULT NULL,
  "is_public" BOOLEAN DEFAULT true,
  "is_user" BOOLEAN NOT NULL,
  "created_at" DATETIME DEFAULT (CURRENT_TIMESTAMP) NOT NULL,
  CHECK (
    is_user = 0
    OR email IS NOT NULL
  ),
    CHECK (
    is_user = 0
    OR first_name IS NOT NULL
  ),
      CHECK (
    is_user = 0
    OR last_name IS NOT NULL
  ),
      CHECK (
    is_user = 0
    OR gender IS NOT NULL
  )
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_profile_email ON profile(email);

CREATE UNIQUE INDEX IF NOT EXISTS idx_profile_display_name ON profile(display_name);

CREATE INDEX IF NOT EXISTS idx_profile_created_at ON profile(created_at);

CREATE INDEX IF NOT EXISTS idx_profile_is_public ON profile(is_public);

CREATE INDEX IF NOT EXISTS idx_profile_is_user ON profile(is_user);