package app

import (
	"log"

	"github.com/fabianherreracruz/bre-b-pse-app/internal/config"
	"github.com/fabianherreracruz/bre-b-pse-app/internal/db"
	"github.com/fabianherreracruz/bre-b-pse-app/internal/handlers"
	"github.com/fabianherreracruz/bre-b-pse-app/internal/middleware"
	"github.com/fabianherreracruz/bre-b-pse-app/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	Router *gin.Engine
	DB     *gorm.DB
	Config *config.Config
}

func New(cfg *config.Config) (*App, error) {
	// Inicializar base de datos
	database, err := db.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Ejecutar migraciones
	if err := db.MigrateDatabase(database); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Inicializar router
	router := gin.Default()

	// Inicializar servicios
	epaycoService := services.NewEPaycoService(cfg)
	notifyService := services.NewNotificationService(cfg)
	splitService := services.NewSplitService(database)
	reportService := services.NewReportService(database)

	// Inicializar handlers
	paymentHandler := handlers.NewPaymentHandler(database, epaycoService, notifyService, splitService)
	reportHandler := handlers.NewReportHandler(database, reportService)

	// Configurar rutas
	setupRoutes(router, paymentHandler, reportHandler, cfg)

	return &App{
		Router: router,
		DB:     database,
		Config: cfg,
	}, nil
}

func setupRoutes(router *gin.Engine, paymentHandler *handlers.PaymentHandler, reportHandler *handlers.ReportHandler, cfg *config.Config) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1
	v1 := router.Group("/api/v1")

	// Payments
	payments := v1.Group("/payments")
	{
		payments.POST("/create", paymentHandler.CreatePayment)
		payments.GET("/verify/:reference", paymentHandler.VerifyPayment)
		payments.GET("/status/:id", paymentHandler.GetPaymentStatus)
		payments.POST("/webhook", paymentHandler.WebhookPayment)
	}

	// Reports (protegidas por middleware)
	reports := v1.Group("/reports")
	reports.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		reports.GET("/by-date", reportHandler.GetReportByDateRange)
		reports.GET("/export-excel", reportHandler.ExportReportToExcel)
		reports.GET("/statistics", reportHandler.GetStatistics)
		reports.GET("/monthly", reportHandler.GetMonthlyReport)
		reports.GET("/yearly", reportHandler.GetYearlyReport)
	}

	log.Println("✅ Routes configured")
}

func (a *App) Run(port string) error {
	return a.Router.Run(":" + port)
}
