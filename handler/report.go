package handler

import (
	"app-inventory/dto"
	"app-inventory/service"
	"app-inventory/utils"
	"net/http"

	"go.uber.org/zap"
)

type ReportHandler struct {
	service service.ReportServiceInterface
	logger  *zap.Logger
}

// constructor
func NewReportHandler(service service.ReportServiceInterface, log *zap.Logger) ReportHandler {
	return ReportHandler{
		service: service,
		logger:  log,
	}
}

func (h *ReportHandler) Report(w http.ResponseWriter, r *http.Request) {

	report, err := h.service.Report()
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusInternalServerError, "Failed to fetch Report: "+err.Error(), nil)
		return
	}
	response := dto.ReportResponse{
		TotalTransactions: report.TotalTransactions,
		TotalItemsSold:    report.TotalItemsSold,
		TotalRevenue:      report.TotalRevenue,
	}

	utils.ResponseSuccess(w, http.StatusOK, "success get data", response)

}
