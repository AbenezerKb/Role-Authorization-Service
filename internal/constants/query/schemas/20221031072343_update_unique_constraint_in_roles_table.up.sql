    create unique index roles_name_tenant_id_deleted_at on roles(name,tenant_id,deleted_at) where deleted_at IS NULL;
