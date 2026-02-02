-- name: GetTatkaData :one
SELECT * from 
tatkal_config
where train_id = $1;