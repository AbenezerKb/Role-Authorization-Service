CREATE TABLE "permissions" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "status" boolean NOT NULL DEFAULT true,
    "name" varchar NOT NULL,
    "description" varchar NOT NULL,
    "statment" JSON NOT NULL,
    "created_at" timestamptz NOT NULL default now(),
    "updated_at" timestamptz NOT NULL default now()
);

CREATE TABLE "permission_domains" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "permission_id" UUID NOT NULL,
    "domain_id" UUID NOT NULL,
    "created_at" timestamptz NOT NULL default now(),
    "updated_at" timestamptz NOT NULL default now()
);

ALTER TABLE "permission_domains" ADD FOREIGN KEY ("permission_id") REFERENCES "permissions" ("id") ON DELETE CASCADE;

ALTER TABLE "permission_domains" ADD FOREIGN KEY ("domain_id") REFERENCES "domains" ("id") ON DELETE CASCADE;

CREATE TABLE "permissions_hierarchy" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "parent" UUID NOT NULL,
    "child" UUID NOT NULL,
    "created_at" timestamptz NOT NULL default now(),
    "updated_at" timestamptz NOT NULL default now()
);

ALTER TABLE "permissions_hierarchy" ADD FOREIGN KEY ("parent") REFERENCES "permissions" ("id");

ALTER TABLE "permissions_hierarchy" ADD FOREIGN KEY ("child") REFERENCES "permissions" ("id");
