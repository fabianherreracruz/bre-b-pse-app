package services

import (
	"fmt"
	"log"
	"time"

	"github.com/fabianherreracruz/bre-b-pse-app/internal/models"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type ReportService struct {
	db *gorm.DB
}

type ReportData struct {
	TotalRecaudos     int64
	TotalAmount       float64
	SuccessfulCount   int64
	FailedCount       int64
	PendingCount      int64
	Recaudos          []models.Recaudo
	Splits            []models.Split
	AverageAmount     float64
}

func NewReportService(db *gorm.DB) *ReportService {
	return &ReportService{db: db}
}

// ObtenerReportePorFecha obtiene reporte de un rango de fechas
func (r *ReportService) ObtenerReportePorFecha(startDate, endDate time.Time, userID *uuid.UUID) (*ReportData, error) {
	var data ReportData
	
	query := r.db.Where("created_at >= ? AND created_at <= ?", startDate, endDate)
	if userID != nil {
		query = query.Where("recaudador_id = ?", userID)
	}

	// Total recaudos
	if err := query.Model(&models.Recaudo{}).Count(&data.TotalRecaudos).Error; err != nil {
		log.Printf("Error counting recaudos: %v", err)
		return nil, err
	}

	// Total monto
	if err := query.Model(&models.Recaudo{}).Select("COALESCE(SUM(amount), 0)").Row().Scan(&data.TotalAmount); err != nil {
		log.Printf("Error calculating total: %v", err)
		return nil, err
	}

	// Recaudos por estado
	query.Model(&models.Recaudo{}).Where("status = ?", "completado").Count(&data.SuccessfulCount)
	query.Model(&models.Recaudo{}).Where("status = ?", "fallido").Count(&data.FailedCount)
	query.Model(&models.Recaudo{}).Where("status = ?", "pendiente").Count(&data.PendingCount)

	// Promedio
	if data.TotalRecaudos > 0 {
		data.AverageAmount = data.TotalAmount / float64(data.TotalRecaudos)
	}

	// Obtener recaudos detallados
	if err := query.Find(&data.Recaudos).Error; err != nil {
		log.Printf("Error fetching recaudos: %v", err)
		return nil, err
	}

	// Obtener splits
	splitQuery := r.db.Joins("JOIN recaudos ON splits.recaudo_id = recaudos.id").
		Where("recaudos.created_at >= ? AND recaudos.created_at <= ?", startDate, endDate)
	if userID != nil {
		splitQuery = splitQuery.Where("recaudos.recaudador_id = ?", userID)
	}

	if err := splitQuery.Find(&data.Splits).Error; err != nil {
		log.Printf("Error fetching splits: %v", err)
	}

	return &data, nil
}

// GenerarReporteExcel genera un reporte en Excel
func (r *ReportService) GenerarReporteExcel(data *ReportData, filename string) error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Error closing Excel file: %v", err)
		}
	}()

	// Hoja 1: Resumen
	f.SetCellValue("Sheet1", "A1", "REPORTE DE RECAUDOS BRE-B PSE")
	f.SetCellValue("Sheet1", "A2", fmt.Sprintf("Total Recaudos: %d", data.TotalRecaudos))
	f.SetCellValue("Sheet1", "A3", fmt.Sprintf("Monto Total: $%.2f", data.TotalAmount))
	f.SetCellValue("Sheet1", "A4", fmt.Sprintf("Exitosos: %d", data.SuccessfulCount))
	f.SetCellValue("Sheet1", "A5", fmt.Sprintf("Fallidos: %d", data.FailedCount))
	f.SetCellValue("Sheet1", "A6", fmt.Sprintf("Pendientes: %d", data.PendingCount))
	f.SetCellValue("Sheet1", "A7", fmt.Sprintf("Promedio: $%.2f", data.AverageAmount))

	// Hoja 2: Detalles de Recaudos
	f.NewSheet("Recaudos")
	headers := []string{"ID", "Referencia", "Monto", "Estado", "Correo Cliente", "Teléfono", "Fecha", "Actualizado"}
	for i, header := range headers {
		f.SetCellValue("Recaudos", fmt.Sprintf("%c1", 'A'+i), header)
	}

	for idx, recaudo := range data.Recaudos {
		row := idx + 2
		f.SetCellValue("Recaudos", fmt.Sprintf("A%d", row), recaudo.ID.String())
		f.SetCellValue("Recaudos", fmt.Sprintf("B%d", row), recaudo.ReferenceCode)
		f.SetCellValue("Recaudos", fmt.Sprintf("C%d", row), recaudo.Amount)
		f.SetCellValue("Recaudos", fmt.Sprintf("D%d", row), recaudo.Status)
		f.SetCellValue("Recaudos", fmt.Sprintf("E%d", row), recaudo.CustomerEmail)
		f.SetCellValue("Recaudos", fmt.Sprintf("F%d", row), recaudo.CustomerPhone)
		f.SetCellValue("Recaudos", fmt.Sprintf("G%d", row), recaudo.CreatedAt)
		f.SetCellValue("Recaudos", fmt.Sprintf("H%d", row), recaudo.UpdatedAt)
	}

	// Hoja 3: Splits
	f.NewSheet("Splits")
	splitHeaders := []string{"ID", "Recaudo ID", "Beneficiario ID", "Porcentaje", "Monto", "Estado", "Pagado En"}
	for i, header := range splitHeaders {
		f.SetCellValue("Splits", fmt.Sprintf("%c1", 'A'+i), header)
	}

	for idx, split := range data.Splits {
		row := idx + 2
		f.SetCellValue("Splits", fmt.Sprintf("A%d", row), split.ID.String())
		f.SetCellValue("Splits", fmt.Sprintf("B%d", row), split.RecaudoID.String())
		f.SetCellValue("Splits", fmt.Sprintf("C%d", row), split.BeneficiarioID.String())
		f.SetCellValue("Splits", fmt.Sprintf("D%d", row), split.Percentage)
		f.SetCellValue("Splits", fmt.Sprintf("E%d", row), split.Amount)
		f.SetCellValue("Splits", fmt.Sprintf("F%d", row), split.Status)
		if split.PaidAt.Valid {
			f.SetCellValue("Splits", fmt.Sprintf("G%d", row), split.PaidAt.Time)
		}
	}

	if err := f.SaveAs(filename); err != nil {
		log.Printf("Error saving Excel file: %v", err)
		return err
	}

	log.Printf("✅ Excel report generated: %s", filename)
	return nil
}

// ObtenerEstadísticas obtiene estadísticas por período
func (r *ReportService) ObtenerEstadísticas(userID *uuid.UUID) (map[string]interface{}, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())

	stats := make(map[string]interface{})

	// Mensual
	monthlyData, err := r.ObtenerReportePorFecha(startOfMonth, now, userID)
	if err != nil {
		return nil, err
	}
	stats["monthly"] = monthlyData

	// Anual
	yearlyData, err := r.ObtenerReportePorFecha(startOfYear, now, userID)
	if err != nil {
		return nil, err
	}
	stats["yearly"] = yearlyData

	// Últimos 7 días
	weekData, err := r.ObtenerReportePorFecha(now.AddDate(0, 0, -7), now, userID)
	if err != nil {
		return nil, err
	}
	stats["weekly"] = weekData

	return stats, nil
}
