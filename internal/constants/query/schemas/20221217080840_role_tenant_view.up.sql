create view role_tenant as 
    select r.name as name ,
    r.created_at as created_at,
     r.id as id,
      r.status as status,
      r.updated_at as updated_at,
      t.tenant_name as tenant_name,
      t.service_id as service_id
    from roles r join tenants t on r.tenant_id=t.id
    where r.deleted_at IS NULL;