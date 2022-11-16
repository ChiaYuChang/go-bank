-- name: DoTransfer :exec
  WITH _data (id, amount) AS (
VALUES
    (sqlc.arg(src_id)::bigint, -sqlc.arg(amount)::decimal),
    (sqlc.arg(dst_id)::bigint,  sqlc.arg(amount)::decimal)
)
UPDATE accounts AS a
   SET balance = balance + _data.amount
  FROM _data
 WHERE a.id = _data.id;

-- name: CreateTransferRecord :one
INSERT INTO transfers (
  src_id, dst_id, amount, status
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateTransferStatus :one
UPDATE transfers SET
    status = $2
WHERE id = $1
RETURNING *;