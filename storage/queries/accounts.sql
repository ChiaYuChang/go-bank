-- name: CreateAccount :one
INSERT INTO accounts (
    owner, balance, currency
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- -- name: BulkCreateAccounts :copyfrom
-- INSERT INTO accounts (
--     owner, balance, currency
-- ) VALUES (
--     $1, $2, $3   
-- );

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateAccount :one
UPDATE accounts SET
    owner = $1,
    balance = $2,
    currency = $3,
    updated_at = $4
WHERE id = $5 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateAccountBalance :one
UPDATE accounts SET
    balance = $1,
    updated_at = $2
WHERE id = $3 AND deleted_at IS NULL
RETURNING *;  

-- name: AccountBalanceDeposit :one
UPDATE accounts SET
    balance = balance + sqlc.arg(amount),
    updated_at = sqlc.arg(updated_at)
WHERE id = sqlc.arg(id) AND deleted_at IS NULL
RETURNING *; 

-- name: AccountBalanceWithdraw :one
UPDATE accounts SET
    balance = balance - sqlc.arg(amount),
    updated_at = sqlc.arg(updated_at)
WHERE id = sqlc.arg(id) AND deleted_at IS NULL
RETURNING *; 

-- name: ListAccounts :many
SELECT * FROM accounts
WHERE deleted_at IS NULL
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: DeleteAccount :exec
UPDATE accounts SET
    deleted_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: HardDeleteAccount :exec
DELETE FROM accounts WHERE id = $1;