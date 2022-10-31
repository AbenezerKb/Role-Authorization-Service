CREATE TABLE "domains" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" varchar NOT NULL ,
    "status" status NOT NULL DEFAULT 'ACTIVE',
    "deleted_at" timestamptz,
    "service_id" UUID NOT NULL,
    "created_at" timestamptz NOT NULL default now(),
    "updated_at" timestamptz NOT NULL default now(),
    UNIQUE("name","service_id","deleted_at") where deleted_at IS NULL
);

ALTER TABLE "domains" ADD FOREIGN KEY ("service_id") REFERENCES "services" ("id") ON DELETE CASCADE;
