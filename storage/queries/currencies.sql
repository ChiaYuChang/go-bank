-- name: CreateCurrency :one
INSERT INTO currencies (
    name, abbr
) VALUES (
    $1, $2
)
RETURNING *;

-- -- name: BulkCreateCurrencies :copyfrom
-- INSERT INTO currencies (
--     name, abbr
-- ) VALUES (
--     $1, $2
-- );

-- name: GetCurrency :one
SELECT * FROM currencies
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateCurrency :one
UPDATE currencies SET
    name = $1,
    abbr = $2
WHERE id = $3 AND deleted_at IS NULL
RETURNING *;

-- name: ListCurrencies :many
SELECT * FROM currencies
WHERE deleted_at IS NULL
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: DeleteCurrency :exec
UPDATE currencies SET
    deleted_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: HardDeleteCurrency :exec
DELETE FROM currencies WHERE id = $1;