-- name: CreateTenent :exec 
INSERT INTO tenants (
domain_id,
tenant_name,
service_id

) VALUES (
 $1,$2,$3
) ;

-- name: GetTenentWithNameAndServiceId :one 
SELECT * FROM tenants WHERE 
tenant_name = $1 AND service_id = $2 AND deleted_at IS NULL;

