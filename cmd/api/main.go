package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "lqstudio-backend/docs" // Import docs for Swagger
	"lqstudio-backend/internal/config"
	"lqstudio-backend/internal/database"
	"lqstudio-backend/internal/database/sqlc"
	"lqstudio-backend/internal/handlers"
	"lqstudio-backend/internal/repositories"
	"lqstudio-backend/internal/services"
	"lqstudio-backend/pkg/email"
	"lqstudio-backend/pkg/logger"
	"lqstudio-backend/pkg/metrics"

	"go.uber.org/zap"
)

// @title LQ Studio Photography Booking API
// @version 1.0
// @description API for photography studio booking system
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@lqstudio.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// Initialize logger
	if err := logger.Init(cfg.Server.Env); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	log := logger.Log
	log.Info("Starting LQ Studio Photography Booking System",
		zap.String("env", cfg.Server.Env),
		zap.Int("port", cfg.Server.Port),
	)
	log.Info("CORS Configuration",
		zap.Strings("allowed_origins", cfg.CORS.AllowedOrigins),
	)

	// Connect to database
	log.Info("Connecting to database...")
	dbConn, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer dbConn.Close()
	log.Info("Database connected successfully")

	// Run database migrations automatically
	log.Info("Running database migrations...")
	sqlDB, err := dbConn.GetStdDB()
	if err != nil {
		log.Fatal("Failed to get database/sql connection for migrations", zap.Error(err))
	}
	defer sqlDB.Close()

	// Migrations directory relative to the binary location
	migrationsPath := "./migrations"
	if err := database.RunMigrations(sqlDB, migrationsPath); err != nil {
		log.Fatal("Failed to run database migrations", zap.Error(err))
	}

	// Get and log current migration version
	version, err := database.GetMigrationStatus(sqlDB)
	if err != nil {
		log.Warn("Failed to get migration version", zap.Error(err))
	} else {
		log.Info("Database migrations completed successfully", zap.Int64("version", version))
	}

	// Initialize Prometheus metrics
	log.Info("Initializing Prometheus metrics...")
	metrics.Init(dbConn.Pool)
	log.Info("Prometheus metrics initialized successfully")

	// Initialize sqlc queries with pgxpool
	queries := sqlc.New(dbConn.Pool)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(queries)
	themeRepo := repositories.NewThemeRepository(queries)
	packageRepo := repositories.NewPackageRepository(queries)
	addonRepo := repositories.NewAddonRepository(queries)
	bookingRepo := repositories.NewBookingRepository(dbConn.Pool, queries)

	// Initialize email client (optional - can be nil if not configured)
	var emailClient *email.Client
	if cfg.Email.APIKey != "" && cfg.Email.From != "" {
		var err error
		emailClient, err = email.NewClient(
			cfg.Email.APIKey,
			cfg.Email.From,
			cfg.Email.AdminTo,
			log,
		)
		if err != nil {
			log.Warn("Failed to initialize email client, email notifications will be disabled",
				zap.Error(err),
			)
		} else {
			log.Info("Email client initialized successfully",
				zap.String("from", cfg.Email.From),
				zap.String("admin_to", cfg.Email.AdminTo),
			)
		}
	} else {
		log.Info("Email configuration not set, email notifications will be disabled")
	}

	// Initialize services
	packageSvc := services.NewPackageService(packageRepo)
	themeSvc := services.NewThemeService(themeRepo)
	addonSvc := services.NewAddonService(addonRepo)
	bookingService := services.NewBookingService(packageRepo, themeRepo, addonRepo, bookingRepo, emailClient, log)

	// Initialize HTTP handlers
	healthHandler := handlers.NewHealthHandler(dbConn)
	availabilityHandler := handlers.NewAvailabilityHandler(bookingService)
	bookingHandler := handlers.NewBookingHandler(bookingService, &cfg.Upload)
	adminHandler := handlers.NewAdminHandler(packageSvc, themeSvc, addonSvc, bookingService)
	authHandler := handlers.NewAuthHandler(userRepo, cfg, log)

	// Setup router
	router := handlers.NewRouter(
		log,
		cfg.CORS.AllowedOrigins,
		cfg.JWT.Secret,
		healthHandler,
		availabilityHandler,
		bookingHandler,
		adminHandler,
		authHandler,
	)
	e := router.Setup()

	// Start server
	go func() {
		addr := fmt.Sprintf(":%d", cfg.Server.Port)
		log.Info("Server starting", zap.String("address", addr))
		if err := e.Start(addr); err != nil {
			log.Info("Server shutting down")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	// Graceful shutdown with 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	// Shutdown metrics collection
	metrics.Shutdown()

	log.Info("Server shutdown complete")
}
