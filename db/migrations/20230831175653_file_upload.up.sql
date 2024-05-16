CREATE TABLE IF NOT EXISTS "images"(
    "id" SERIAL PRIMARY KEY,
    "image_url" TEXT NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())

)