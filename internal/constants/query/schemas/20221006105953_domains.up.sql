CREATE TABLE "domains" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" varchar NOT NULL UNIQUE,
    "deleted_at" timestamptz,
    "service_id" UUID NOT NULL,
    "created_at" timestamptz NOT NULL default now(),
    "updated_at" timestamptz NOT NULL default now()
);

ALTER TABLE "domains" ADD FOREIGN KEY ("service_id") REFERENCES "services" ("id") ON DELETE CASCADE;
