
-- name: GetRecordByID :one
SELECT * FROM record WHERE id = $1;

-- name: GetRecords :many
SELECT * FROM record;

-- name: GetRecordByDomainName :one
SELECT * FROM record WHERE domain_name = $1;
