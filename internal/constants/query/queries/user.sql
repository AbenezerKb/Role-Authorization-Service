-- name: RegisterUser :exec 
INSERT INTO users (
user_id,
service_id
) VALUES (
 $1,$2
);

-- name: GetUserWithUserIdAndServiceId :one 
SELECT * FROM users WHERE 
user_id = $1 AND service_id = $2 AND deleted_at IS NULL;

-- name: UpdateUserStatus :one
UPDATE users SET status = $1 WHERE user_id = $2 AND service_id=$3 AND deleted_at IS NULL RETURNING id;