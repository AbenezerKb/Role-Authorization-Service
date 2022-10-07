CREATE TABLE "users" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "user_id" UUID NOT NULL UNIQUE,
    "status" boolean NOT NULL DEFAULT true,
    "deleted_at" timestamptz,
    "created_at" timestamptz NOT NULL default now(),
    "updated_at" timestamptz NOT NULL default now()
);

CREATE TABLE "tenant_users_roles" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "tenant_id" varchar NOT NULL,
    "user_id" UUID NOT NULL,
    "role_id" UUID NOT NULL,
    "status" boolean NOT NULL default true,
    "deleted_at" timestamptz,
    "created_at" timestamptz NOT NULL default now(),
    "updated_at" timestamptz NOT NULL default now()
);
ALTER TABLE "tenant_users_roles" ADD FOREIGN KEY ("tenant_id") REFERENCES "tenants" ("tenant_name");

ALTER TABLE "tenant_users_roles" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");

ALTER TABLE "tenant_users_roles" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");