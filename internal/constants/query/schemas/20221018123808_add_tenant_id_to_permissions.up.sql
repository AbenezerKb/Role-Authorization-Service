ALTER TABLE permissions ADD "tenant_id" UUID;
ALTER TABLE permissions ADD FOREIGN KEY("tenant_id") REFERENCES tenants("id") ON DELETE CASCADE;
