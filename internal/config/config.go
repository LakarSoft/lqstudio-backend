package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	CORS     CORSConfig
	Email    EmailConfig
	Upload   UploadConfig
}

type ServerConfig struct {
	Port int
	Env  string
}

type DatabaseConfig struct {
	Host            string
	Port            string
	Database        string
	Username        string
	Password        string
	Schema          string
	Timezone        string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

type JWTConfig struct {
	Secret      string
	ExpiryHours int
}

type CORSConfig struct {
	AllowedOrigins []string
}

type EmailConfig struct {
	APIKey  string // Resend API key
	From    string // Sender email
	AdminTo string // Admin notification email
}

type UploadConfig struct {
	MaxFileSize  int64    // Max file size in bytes
	AllowedTypes []string // Allowed MIME types
	StoragePath  string   // Where to save uploaded files
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid PORT: %w", err)
	}

	jwtExpiryHours, err := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRY_HOURS: %w", err)
	}

	// Database pool configuration
	maxConns := int32(getEnvInt("DB_MAX_CONNS", 25))
	minConns := int32(getEnvInt("DB_MIN_CONNS", 5))
	maxConnLifetime := time.Duration(getEnvInt("DB_MAX_CONN_LIFETIME_MINUTES", 60)) * time.Minute
	maxConnIdleTime := time.Duration(getEnvInt("DB_MAX_CONN_IDLE_TIME_MINUTES", 10)) * time.Minute

	// Upload configuration
	maxFileSize := int64(getEnvInt("UPLOAD_MAX_FILE_SIZE", 5242880)) // Default 5MB
	allowedTypes := strings.Split(getEnv("UPLOAD_ALLOWED_TYPES", "image/jpeg,image/png,image/jpg"), ",")
	storagePath := getEnv("UPLOAD_STORAGE_PATH", "./uploads/payment-screenshots")

	// Parse CORS origins and trim whitespace
	corsOriginsRaw := strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "*"), ",")
	corsOrigins := make([]string, 0, len(corsOriginsRaw))
	for _, origin := range corsOriginsRaw {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			corsOrigins = append(corsOrigins, trimmed)
		}
	}

	// Validate required configurations
	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required but not set")
	}

	emailAPIKey := getEnv("EMAIL_API_KEY", "")
	emailFrom := getEnv("EMAIL_FROM", "")
	emailAdminTo := getEnv("EMAIL_ADMIN_TO", "")

	// Email validation - warn if not set but don't fail
	if emailAPIKey == "" {
		fmt.Println("WARNING: EMAIL_API_KEY not set - email notifications will be disabled")
	}
	if emailFrom == "" {
		fmt.Println("WARNING: EMAIL_FROM not set - email notifications will be disabled")
	}
	if emailAdminTo == "" {
		fmt.Println("WARNING: EMAIL_ADMIN_TO not set - admin notifications will be disabled")
	}

	return &Config{
		Server: ServerConfig{
			Port: port,
			Env:  getEnv("APP_ENV", "local"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			Database:        getEnv("DB_DATABASE", "lqstudio"),
			Username:        getEnv("DB_USERNAME", "lqstudio_user"),
			Password:        getEnv("DB_PASSWORD", "lqstudio_password_123"),
			Schema:          getEnv("DB_SCHEMA", "public"),
			Timezone:        getEnv("DB_TIMEZONE", "Asia/Kuala_Lumpur"),
			MaxConns:        maxConns,
			MinConns:        minConns,
			MaxConnLifetime: maxConnLifetime,
			MaxConnIdleTime: maxConnIdleTime,
		},
		JWT: JWTConfig{
			Secret:      jwtSecret,
			ExpiryHours: jwtExpiryHours,
		},
		CORS: CORSConfig{
			AllowedOrigins: corsOrigins,
		},
		Email: EmailConfig{
			APIKey:  emailAPIKey,
			From:    emailFrom,
			AdminTo: emailAdminTo,
		},
		Upload: UploadConfig{
			MaxFileSize:  maxFileSize,
			AllowedTypes: allowedTypes,
			StoragePath:  storagePath,
		},
	}, nil
}

// ConnectionString returns the PostgreSQL connection string
func (c *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s&timezone=%s",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.Schema,
		c.Timezone,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// StudioConfig holds business configuration
type StudioConfig struct {
	OpenHour             int
	CloseHour            int
	SlotDurationMinutes  int
	PaymentQRCodeURL     string
	PaymentBankName      string
	PaymentAccountNumber string
	PaymentAccountName   string
	WhatsAppNumber       string
}

// LoadConfig loads configuration from environment
func LoadConfig() *StudioConfig {
	return &StudioConfig{
		OpenHour:             getEnvInt("STUDIO_OPEN_HOUR", 9),
		CloseHour:            getEnvInt("STUDIO_CLOSE_HOUR", 17),
		SlotDurationMinutes:  getEnvInt("SLOT_DURATION_MINUTES", 20),
		PaymentQRCodeURL:     getEnv("PAYMENT_QR_CODE_URL", ""),
		PaymentBankName:      getEnv("PAYMENT_BANK_NAME", "Maybank"),
		PaymentAccountNumber: getEnv("PAYMENT_ACCOUNT_NUMBER", ""),
		PaymentAccountName:   getEnv("PAYMENT_ACCOUNT_NAME", "LQ Studio"),
		WhatsAppNumber:       getEnv("WHATSAPP_NUMBER", "+60123456789"),
	}
}

// GetAvailableSlots generates all possible time slots for a date
func (c *StudioConfig) GetAvailableSlots(date time.Time) []time.Time {
	slots := []time.Time{}

	// Start at opening hour on the given date
	current := time.Date(date.Year(), date.Month(), date.Day(), c.OpenHour, 0, 0, 0, date.Location())
	endTime := time.Date(date.Year(), date.Month(), date.Day(), c.CloseHour, 0, 0, 0, date.Location())

	for current.Before(endTime) {
		slots = append(slots, current)
		current = current.Add(time.Duration(c.SlotDurationMinutes) * time.Minute)
	}

	return slots
}
