package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"kasir-api/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TransactionPG is a PostgreSQL implementation of TransactionRepository.
type TransactionPG struct {
	pool *pgxpool.Pool
}

// NewTransactionPG creates a new PostgreSQL transaction repository.
func NewTransactionPG(pool *pgxpool.Pool) *TransactionPG {
	return &TransactionPG{pool: pool}
}

// CreateTransaction creates a transaction from checkout items in a single DB transaction:
// for each item loads product (name, price, stock), decrements stock, builds details and totals,
// inserts the transaction and its details, then commits.
func (r *TransactionPG) CreateTransaction(items []domain.CheckoutItem) (*domain.Transaction, error) {
	tx, err := r.pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	totalAmount := 0
	details := make([]domain.TransactionDetail, 0, len(items))

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := tx.QueryRow(context.Background(),
			"SELECT nama, harga, stok FROM products WHERE id = $1", item.ProductID).
			Scan(&productName, &productPrice, &stock)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("product id %d not found", item.ProductID)
			}
			return nil, err
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec(context.Background(),
			"UPDATE products SET stok = stok - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, domain.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	var createdAt time.Time
	err = tx.QueryRow(context.Background(),
		"INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id, created_at", totalAmount).
		Scan(&transactionID, &createdAt)
	if err != nil {
		return nil, err
	}

	for i := range details {
		details[i].TransactionID = transactionID
		err = tx.QueryRow(context.Background(),
			"INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4) RETURNING id",
			transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal).
			Scan(&details[i].ID)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		return nil, err
	}

	return &domain.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		CreatedAt:   createdAt,
		Details:     details,
	}, nil
}
