package handler

import (
	"app-inventory/dto"
	"app-inventory/service"
	"app-inventory/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ReportHandler struct {
	service service.ReportServiceInterface
	logger  *zap.Logger
}

func NewReportHandler(service service.ReportServiceInterface, log *zap.Logger) ReportHandler {
	return ReportHandler{
		service: service,
		logger:  log,
	}
}

func (h *ReportHandler) Report(c *gin.Context) {
	report, err := h.service.Report(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get report on service", zap.Error(err))
		utils.ResponseError(c, http.StatusInternalServerError, "Gagal memuat ringkasan laporan: "+err.Error(), nil)
		return
	}

	response := dto.ReportResponse{
		TotalTransactions: report.TotalTransactions,
		TotalItemsSold:    report.TotalItemsSold,
		TotalRevenue:      report.TotalRevenue,
	}

	utils.ResponseSuccess(c, http.StatusOK, "Berhasil memuat ringkasan laporan", response)
}
