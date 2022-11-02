create unique index domains_name_service_id_key on domains(name,service_id) where deleted_at IS NULL;
    