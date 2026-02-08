package usecase

import (
	"kasir-api/internal/domain"
	"kasir-api/internal/repository"
)

type ReportUsecase struct {
	txRepo repository.TransactionRepository
}

func NewReportUsecase(txRepo repository.TransactionRepository) *ReportUsecase {
	return &ReportUsecase{txRepo: txRepo}
}

func (u *ReportUsecase) SummaryHariIni() (*domain.SummaryHariIni, error) {
	return u.txRepo.GetSummaryHariIni()
}
