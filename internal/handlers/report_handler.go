package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fabianherreracruz/bre-b-pse-app/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReportHandler struct {
	db            *gorm.DB
	reportService *services.ReportService
}

func NewReportHandler(db *gorm.DB, reportService *services.ReportService) *ReportHandler {
	return &ReportHandler{
		db:            db,
		reportService: reportService,
	}
}

// GetReportByDateRange obtiene reporte por rango de fechas
func (h *ReportHandler) GetReportByDateRange(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// Parsear fechas
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
		return
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format"})
		return
	}

	// Obtener user_id del contexto
	userID, exists := c.Get("user_id")
	var uid *uuid.UUID
	if exists {
		id := userID.(uuid.UUID)
		uid = &id
	}

	// Obtener reporte
	data, err := h.reportService.ObtenerReportePorFecha(start, end, uid)
	if err != nil {
		log.Printf("Error fetching report: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch report"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// ExportReportToExcel exporta reporte a Excel
func (h *ReportHandler) ExportReportToExcel(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// Parsear fechas
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
		return
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format"})
		return
	}

	// Obtener user_id del contexto
	userID, exists := c.Get("user_id")
	var uid *uuid.UUID
	if exists {
		id := userID.(uuid.UUID)
		uid = &id
	}

	// Obtener reporte
	data, err := h.reportService.ObtenerReportePorFecha(start, end, uid)
	if err != nil {
		log.Printf("Error fetching report: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch report"})
		return
	}

	// Generar nombre de archivo
	filename := fmt.Sprintf("report_%s_%s.xlsx", start.Format("20060102"), end.Format("20060102"))

	// Exportar a Excel
	if err := h.reportService.GenerarReporteExcel(data, filename); err != nil {
		log.Printf("Error generating Excel: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate report"})
		return
	}

	// Enviar archivo
	c.File(filename)
}

// GetStatistics obtiene estadísticas
func (h *ReportHandler) GetStatistics(c *gin.Context) {
	// Obtener user_id del contexto
	userID, exists := c.Get("user_id")
	var uid *uuid.UUID
	if exists {
		id := userID.(uuid.UUID)
		uid = &id
	}

	// Obtener estadísticas
	stats, err := h.reportService.ObtenerEstadísticas(uid)
	if err != nil {
		log.Printf("Error fetching statistics: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetMonthlyReport obtiene reporte mensual
func (h *ReportHandler) GetMonthlyReport(c *gin.Context) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 1, 0).Add(-time.Second)

	// Obtener user_id del contexto
	userID, exists := c.Get("user_id")
	var uid *uuid.UUID
	if exists {
		id := userID.(uuid.UUID)
		uid = &id
	}

	// Obtener reporte
	data, err := h.reportService.ObtenerReportePorFecha(start, end, uid)
	if err != nil {
		log.Printf("Error fetching monthly report: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch report"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetYearlyReport obtiene reporte anual
func (h *ReportHandler) GetYearlyReport(c *gin.Context) {
	now := time.Now()
	start := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	end := now

	// Obtener user_id del contexto
	userID, exists := c.Get("user_id")
	var uid *uuid.UUID
	if exists {
		id := userID.(uuid.UUID)
		uid = &id
	}

	// Obtener reporte
	data, err := h.reportService.ObtenerReportePorFecha(start, end, uid)
	if err != nil {
		log.Printf("Error fetching yearly report: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch report"})
		return
	}

	c.JSON(http.StatusOK, data)
}
