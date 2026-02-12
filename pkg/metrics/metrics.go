package metrics

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Global metrics - exported so they can be used throughout the application
var (
	// HTTP Metrics - following RED methodology (Rate, Errors, Duration)
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "lqstudio_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "lqstudio_http_request_duration_seconds",
			Help: "HTTP request latency in seconds",
			// Buckets optimized for web APIs (10ms to 10s)
			Buckets: []float64{0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0},
		},
		[]string{"method", "path"},
	)

	HTTPRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "lqstudio_http_requests_in_flight",
			Help: "Current number of HTTP requests being processed",
		},
	)

	// Business Metrics
	BookingsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "lqstudio_bookings_total",
			Help: "Total number of bookings created",
		},
	)

	BookingsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "lqstudio_bookings_active",
			Help: "Current number of active bookings (PENDING, CONFIRMED states)",
		},
	)

	EmailNotificationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "lqstudio_email_notifications_total",
			Help: "Total number of emails sent by type",
		},
		[]string{"type"}, // customer_confirmation, admin_notification
	)

	EmailFailuresTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "lqstudio_email_failures_total",
			Help: "Total number of failed email send attempts by type",
		},
		[]string{"type"},
	)

	// Database Metrics
	DBQueriesTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "lqstudio_db_queries_total",
			Help: "Total number of database queries executed",
		},
	)

	DBQueryDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name: "lqstudio_db_query_duration_seconds",
			Help: "Database query execution time in seconds",
			// Buckets optimized for DB queries (1ms to 5s)
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 5.0},
		},
	)

	DBConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "lqstudio_db_connections_active",
			Help: "Current number of active database connections",
		},
	)

	DBConnectionsIdle = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "lqstudio_db_connections_idle",
			Help: "Current number of idle database connections in pool",
		},
	)

	DBConnectionsMax = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "lqstudio_db_connections_max",
			Help: "Maximum allowed database connections",
		},
	)
)

var (
	dbPool *pgxpool.Pool
	stopCh chan struct{}
)

// Init initializes the metrics system and starts background collectors
func Init(pool *pgxpool.Pool) {
	dbPool = pool

	// Note: Go runtime collectors (goroutines, memory, GC) and process collectors
	// are already registered by default in Prometheus, so we don't need to register them again

	// Set static database pool configuration
	if pool != nil {
		config := pool.Config()
		DBConnectionsMax.Set(float64(config.MaxConns))

		// Start background goroutine to collect database pool stats
		stopCh = make(chan struct{})
		go collectDBPoolStats()
	}
}

// Shutdown stops background metric collectors
func Shutdown() {
	if stopCh != nil {
		close(stopCh)
	}
}

// collectDBPoolStats periodically collects database connection pool statistics
func collectDBPoolStats() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if dbPool != nil {
				stats := dbPool.Stat()
				DBConnectionsActive.Set(float64(stats.AcquiredConns()))
				DBConnectionsIdle.Set(float64(stats.IdleConns()))
			}
		case <-stopCh:
			return
		}
	}
}

// RecordDBQuery records a database query execution with timing
func RecordDBQuery(start time.Time) {
	duration := time.Since(start).Seconds()
	DBQueriesTotal.Inc()
	DBQueryDuration.Observe(duration)
}

// UpdateActiveBookingsCount updates the active bookings gauge
// This should be called periodically or when booking status changes
func UpdateActiveBookingsCount(ctx context.Context, count int) {
	BookingsActive.Set(float64(count))
}
