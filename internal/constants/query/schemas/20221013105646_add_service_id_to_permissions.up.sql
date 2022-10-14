ALTER TABLE permissions ADD "service_id" UUID NOT NULL;
ALTER TABLE permissions ADD FOREIGN KEY("service_id") REFERENCES services("id") ON DELETE CASCADE;
