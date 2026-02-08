package repository

import "kasir-api/internal/domain"

// TransactionRepository defines the interface for transaction data access.
type TransactionRepository interface {
	CreateTransaction(items []domain.CheckoutItem) (*domain.Transaction, error)
	GetSummaryHariIni() (*domain.SummaryHariIni, error)
}