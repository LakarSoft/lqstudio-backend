package handlers

import (
	"lqstudio-backend/pkg/middleware"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
)

type Router struct {
	echo                *echo.Echo
	logger              *zap.Logger
	allowedOrigins      []string
	jwtSecret           string
	healthHandler       *HealthHandler
	availabilityHandler *AvailabilityHandler
	bookingHandler      *BookingHandler
	adminHandler        *AdminHandler
	authHandler         *AuthHandler
}

func NewRouter(
	logger *zap.Logger,
	allowedOrigins []string,
	jwtSecret string,
	healthHandler *HealthHandler,
	availabilityHandler *AvailabilityHandler,
	bookingHandler *BookingHandler,
	adminHandler *AdminHandler,
	authHandler *AuthHandler,
) *Router {
	return &Router{
		echo:                echo.New(),
		logger:              logger,
		allowedOrigins:      allowedOrigins,
		jwtSecret:           jwtSecret,
		healthHandler:       healthHandler,
		availabilityHandler: availabilityHandler,
		bookingHandler:      bookingHandler,
		adminHandler:        adminHandler,
		authHandler:         authHandler,
	}
}

func (r *Router) Setup() *echo.Echo {
	// Middleware
	r.echo.Use(echoMiddleware.Recover())
	r.echo.Use(middleware.Metrics())   // Add metrics tracking
	r.echo.Use(middleware.RequestID()) // Add request ID to all requests
	r.echo.Use(middleware.Logger(r.logger))
	r.echo.Use(middleware.CORS(r.allowedOrigins))

	// Custom error handler
	r.echo.HTTPErrorHandler = middleware.ErrorHandler

	// Health check (public)
	r.echo.GET("/health", r.healthHandler.HealthCheck)

	// Prometheus metrics endpoint (public)
	r.echo.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// Swagger documentation
	r.echo.GET("/swagger/*", echoSwagger.WrapHandler)

	// Serve uploaded files (public - for payment screenshots)
	r.echo.Static("/uploads", "./uploads")

	// API routes (using /api instead of /api/v1 to match frontend expectations)
	api := r.echo.Group("/api")

	// ========================================================================
	// Public routes - Packages, Themes, Add-ons
	// ========================================================================
	api.GET("/packages", r.adminHandler.GetActivePackages)
	api.GET("/themes", r.adminHandler.GetActiveThemes)
	api.GET("/addons", r.adminHandler.GetActiveAddons)

	// Availability routes (public)
	// GET /api/themes/:themeId/available-times?date=YYYY-MM-DD
	api.GET("/themes/:themeId/available-times", r.availabilityHandler.GetAvailability)

	// Booking routes (public - customers don't need to login)
	api.POST("/bookings", r.bookingHandler.CreateBooking)
	api.GET("/bookings/:id", r.bookingHandler.GetBookingByID)
	api.POST("/bookings/:id/payment-screenshot", r.bookingHandler.UploadPaymentScreenshot)

	// ========================================================================
	// Admin routes
	// ========================================================================

	// Admin Login (public - no authentication required)
	api.POST("/admin/login", r.authHandler.Login)

	// Protected admin routes (with JWT authentication)
	admin := api.Group("/admin")
	admin.Use(middleware.RequireAuth(r.jwtSecret))
	admin.Use(middleware.RequireRole("ADMIN"))

	// Admin - Package Management
	admin.GET("/packages", r.adminHandler.GetAllPackages)
	admin.GET("/packages/:id", r.adminHandler.GetPackage)
	admin.POST("/packages", r.adminHandler.CreatePackage)
	admin.PUT("/packages/:id", r.adminHandler.UpdatePackage)
	admin.DELETE("/packages/:id", r.adminHandler.DeletePackage)
	admin.PATCH("/packages/:id/toggle-active", r.adminHandler.TogglePackageActive)
	admin.POST("/packages/:id/image", r.adminHandler.UploadPackageImage)

	// Admin - Theme Management
	admin.GET("/themes", r.adminHandler.GetAllThemes)
	admin.GET("/themes/:id", r.adminHandler.GetTheme)
	admin.POST("/themes", r.adminHandler.CreateTheme)
	admin.PUT("/themes/:id", r.adminHandler.UpdateTheme)
	admin.DELETE("/themes/:id", r.adminHandler.DeleteTheme)
	admin.PATCH("/themes/:id/toggle-active", r.adminHandler.ToggleThemeActive)
	admin.POST("/themes/:id/image", r.adminHandler.UploadThemeImage)

	// Admin - Add-on Management
	admin.GET("/addons", r.adminHandler.GetAllAddons)
	admin.GET("/addons/:id", r.adminHandler.GetAddon)
	admin.POST("/addons", r.adminHandler.CreateAddon)
	admin.PUT("/addons/:id", r.adminHandler.UpdateAddon)
	admin.DELETE("/addons/:id", r.adminHandler.DeleteAddon)
	admin.PATCH("/addons/:id/toggle-active", r.adminHandler.ToggleAddonActive)

	// Admin - Booking Management
	// GET /api/admin/bookings?status=PENDING&limit=20&offset=0
	admin.GET("/bookings", r.adminHandler.ListBookings)
	admin.GET("/bookings/:id", r.adminHandler.GetBooking)
	// PATCH /api/admin/bookings/:id/status
	admin.PATCH("/bookings/:id/status", r.adminHandler.UpdateBookingStatus)
	admin.PUT("/bookings/:id", r.bookingHandler.UpdateBooking)

	return r.echo
}
