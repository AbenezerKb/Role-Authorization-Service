CREATE TABLE "roles" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "status" boolean NOT NULL DEFAULT true,
    "name" varchar NOT NULL,
    "tenant_id" UUID NOT NULL,
    "deleted_at" timestamptz,
    "created_at" timestamptz NOT NULL default now(),
    "updated_at" timestamptz NOT NULL default now()
);

ALTER TABLE "roles" add FOREIGN KEY ("tenant_id") REFERENCES "tenants"  ("id") ON DELETE CASCADE;

CREATE TABLE "role_permissions"(
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "permission_id" UUID NOT NULL,
    "role_id" UUID NOT NULL,
    "created_at" timestamptz NOT NULL default now(),
    "updated_at" timestamptz NOT NULL default now()
);

ALTER TABLE "role_permissions" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id") ON DELETE CASCADE;

ALTER TABLE "role_permissions" ADD FOREIGN KEY ("permission_id") REFERENCES "permissions" ("id") ON DELETE CASCADE;
