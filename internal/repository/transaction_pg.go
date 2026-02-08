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

// GetSummaryHariIni returns today's sales summary: total revenue, transaction count, and best-selling product.
func (r *TransactionPG) GetSummaryHariIni() (*domain.SummaryHariIni, error) {
	ctx := context.Background()
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var totalRevenue int
	var totalTransaksi int
	err := r.pool.QueryRow(ctx,
		"SELECT COALESCE(SUM(total_amount), 0), COUNT(*) FROM transactions WHERE created_at >= $1 AND created_at < $2",
		startOfDay, endOfDay).
		Scan(&totalRevenue, &totalTransaksi)
	if err != nil {
		return nil, err
	}

	out := &domain.SummaryHariIni{
		TotalRevenue:   totalRevenue,
		TotalTransaksi: totalTransaksi,
		ProdukTerlaris: domain.ProdukTerlaris{},
	}

	var nama string
	var qtyTerjual int
	err = r.pool.QueryRow(ctx,
		`SELECT p.nama, COALESCE(SUM(td.quantity), 0)
		 FROM transaction_details td
		 JOIN transactions t ON t.id = td.transaction_id
		 JOIN products p ON p.id = td.product_id
		 WHERE t.created_at >= $1 AND t.created_at < $2
		 GROUP BY td.product_id, p.nama
		 ORDER BY SUM(td.quantity) DESC
		 LIMIT 1`,
		startOfDay, endOfDay).
		Scan(&nama, &qtyTerjual)
	if err == nil {
		out.ProdukTerlaris = domain.ProdukTerlaris{Nama: nama, QtyTerjual: qtyTerjual}
	}
	// ErrNoRows means no sales today; leave ProdukTerlaris as zero value
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	return out, nil
}
