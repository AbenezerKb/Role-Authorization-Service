CREATE TABLE "tenants" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "status" boolean NOT NULL DEFAULT true,
    "tenant_name" varchar UNIQUE NOT NULL,
    "service_id" UUID NOT NULL,
    "deleted_at" timestamptz,
    "created_at" timestamptz NOT NULL default now(),
    "updated_at" timestamptz NOT NULL default now()
);

ALTER TABLE "tenants" add FOREIGN KEY ("service_id") REFERENCES "services"  ("id");