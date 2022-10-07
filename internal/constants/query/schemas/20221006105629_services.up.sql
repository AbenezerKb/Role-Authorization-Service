CREATE TABLE "services" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "status" boolean NOT NULL DEFAULT false,
    "name" varchar UNIQUE NOT NULL,
    "password" varchar NOT NULL,
    "deleted_at" timestamptz,
    "created_at" timestamptz NOT NULL default now(),
    "updated_at" timestamptz NOT NULL default now()
);