package usecase

import (
	"kasir-api/internal/domain"
	"kasir-api/internal/repository"
)

type TransactionUsecase struct {
	repo repository.TransactionRepository
}

func NewTransactionUsecase(repo repository.TransactionRepository) *TransactionUsecase {
	return &TransactionUsecase{repo: repo}
}

func (u *TransactionUsecase) Checkout(items []domain.CheckoutItem, useLock bool) (*domain.Transaction, error) {
	return u.repo.CreateTransaction(items)
}
