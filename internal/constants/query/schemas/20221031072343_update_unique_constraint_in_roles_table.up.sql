    create unique index roles_name_tenant_id_key on roles(name,tenant_id) where deleted_at IS NULL;
