// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: currencies.sql

package models

import (
	"context"
)

const createCurrency = `-- name: CreateCurrency :one
INSERT INTO currencies (
    name, abbr
) VALUES (
    $1, $2
)
RETURNING id, name, abbr, created_at, updated_at, deleted_at
`

type CreateCurrencyParams struct {
	Name string `json:"name"`
	Abbr string `json:"abbr"`
}

func (q *Queries) CreateCurrency(ctx context.Context, arg CreateCurrencyParams) (Currency, error) {
	row := q.db.QueryRowContext(ctx, createCurrency, arg.Name, arg.Abbr)
	var i Currency
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Abbr,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const deleteCurrency = `-- name: DeleteCurrency :exec
UPDATE currencies SET
    deleted_at = now()
WHERE id = $1 AND deleted_at IS NULL
`

func (q *Queries) DeleteCurrency(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteCurrency, id)
	return err
}

const getCurrency = `-- name: GetCurrency :one

SELECT id, name, abbr, created_at, updated_at, deleted_at FROM currencies
WHERE id = $1 AND deleted_at IS NULL
`

// -- name: BulkCreateCurrencies :copyfrom
// INSERT INTO currencies (
//
//	name, abbr
//
// ) VALUES (
//
//	$1, $2
//
// );
func (q *Queries) GetCurrency(ctx context.Context, id int32) (Currency, error) {
	row := q.db.QueryRowContext(ctx, getCurrency, id)
	var i Currency
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Abbr,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const hardDeleteCurrency = `-- name: HardDeleteCurrency :exec
DELETE FROM currencies WHERE id = $1
`

func (q *Queries) HardDeleteCurrency(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, hardDeleteCurrency, id)
	return err
}

const listCurrencies = `-- name: ListCurrencies :many
SELECT id, name, abbr, created_at, updated_at, deleted_at FROM currencies
WHERE deleted_at IS NULL
ORDER BY id
LIMIT $1 OFFSET $2
`

type ListCurrenciesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListCurrencies(ctx context.Context, arg ListCurrenciesParams) ([]Currency, error) {
	rows, err := q.db.QueryContext(ctx, listCurrencies, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Currency
	for rows.Next() {
		var i Currency
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Abbr,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateCurrency = `-- name: UpdateCurrency :one
UPDATE currencies SET
    name = $1,
    abbr = $2
WHERE id = $3 AND deleted_at IS NULL
RETURNING id, name, abbr, created_at, updated_at, deleted_at
`

type UpdateCurrencyParams struct {
	Name string `json:"name"`
	Abbr string `json:"abbr"`
	ID   int32  `json:"id"`
}

func (q *Queries) UpdateCurrency(ctx context.Context, arg UpdateCurrencyParams) (Currency, error) {
	row := q.db.QueryRowContext(ctx, updateCurrency, arg.Name, arg.Abbr, arg.ID)
	var i Currency
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Abbr,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}