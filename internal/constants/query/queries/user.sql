-- name: RegisterUser :exec 
INSERT INTO users (
user_id,
service_id
) VALUES (
 $1,$2
);

-- name: GetUserWithUserIdAndServiceId :one 
SELECT * FROM users WHERE 
user_id = $1 AND service_id = $2;
