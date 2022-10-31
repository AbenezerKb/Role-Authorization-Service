create unique index domains_name_service_id_deleted_at_key on domains(name,service_id,deleted_at) where deleted_at IS NULL;
    