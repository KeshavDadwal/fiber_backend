CREATE TABLE IF NOT EXISTS "role"(
	"id" SERIAL PRIMARY KEY,
	"role_name" TEXT NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "permissions"(
	"id" SERIAL PRIMARY KEY,
	"permission_name" TEXT NOT NULL,
    "created_at"timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "role_permission"(
    "id" SERIAL PRIMARY KEY,
    "role_id" INT NOT NULL,
    "permission_id" INT NOT NUll,
    FOREIGN KEY ("role_id") REFERENCES "role" ("id"),
    FOREIGN KEY ("permission_id") REFERENCES "permissions" ("id")
)

