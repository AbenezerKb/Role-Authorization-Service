create unique index users_user_id_service_id_deleted_at_key on users(user_id,service_id,deleted_at)where deleted_at IS NULL;
create unique index tenant_users_roles_user_id_role_id_tenant_id_deleted_at_key on tenant_users_roles (user_id, role_id, tenant_id, deleted_at) WHERE deleted_at IS NULL;
;
