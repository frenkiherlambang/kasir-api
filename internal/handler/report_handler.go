package handler

import (
	"net/http"

	"kasir-api/internal/usecase"
)

type ReportHandler struct {
	uc *usecase.ReportUsecase
}

func NewReportHandler(uc *usecase.ReportUsecase) *ReportHandler {
	return &ReportHandler{uc: uc}
}

// HariIni handles GET /api/report/hari-ini and returns today's sales summary.
func (h *ReportHandler) HariIni(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	sum, err := h.uc.SummaryHariIni()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, sum)
}
