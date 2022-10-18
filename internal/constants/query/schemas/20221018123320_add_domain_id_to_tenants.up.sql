ALTER TABLE tenants ADD "domain_id" UUID NOT NULL;
ALTER TABLE tenants ADD FOREIGN KEY("domain_id") REFERENCES domains("id") ON DELETE CASCADE;
