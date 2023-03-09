// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: payment.sql

package db

import (
	"context"
)

const createPayment = `-- name: CreatePayment :one
INSERT INTO payments (
  merchant_id,
  customer_id,
  amount
) VALUES (
  $1, $2, $3
) RETURNING id, merchant_id, customer_id, amount, created_at
`

type CreatePaymentParams struct {
	MerchantID int64 `json:"merchant_id"`
	CustomerID int64 `json:"customer_id"`
	Amount     int64 `json:"amount"`
}

func (q *Queries) CreatePayment(ctx context.Context, arg CreatePaymentParams) (Payment, error) {
	row := q.db.QueryRowContext(ctx, createPayment, arg.MerchantID, arg.CustomerID, arg.Amount)
	var i Payment
	err := row.Scan(
		&i.ID,
		&i.MerchantID,
		&i.CustomerID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const getPayment = `-- name: GetPayment :one
SELECT id, merchant_id, customer_id, amount, created_at FROM payments
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetPayment(ctx context.Context, id int64) (Payment, error) {
	row := q.db.QueryRowContext(ctx, getPayment, id)
	var i Payment
	err := row.Scan(
		&i.ID,
		&i.MerchantID,
		&i.CustomerID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const listPayments = `-- name: ListPayments :many
SELECT id, merchant_id, customer_id, amount, created_at FROM payments
WHERE 
    merchant_id = $1 OR
    customer_id = $2
ORDER BY id
LIMIT $3
OFFSET $4
`

type ListPaymentsParams struct {
	MerchantID int64 `json:"merchant_id"`
	CustomerID int64 `json:"customer_id"`
	Limit      int32 `json:"limit"`
	Offset     int32 `json:"offset"`
}

func (q *Queries) ListPayments(ctx context.Context, arg ListPaymentsParams) ([]Payment, error) {
	rows, err := q.db.QueryContext(ctx, listPayments,
		arg.MerchantID,
		arg.CustomerID,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Payment{}
	for rows.Next() {
		var i Payment
		if err := rows.Scan(
			&i.ID,
			&i.MerchantID,
			&i.CustomerID,
			&i.Amount,
			&i.CreatedAt,
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
