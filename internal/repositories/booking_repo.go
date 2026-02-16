package repositories

import (
	"context"
	"fmt"
	"time"

	"lqstudio-backend/internal/database/sqlc"
	"lqstudio-backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// BookingRepository handles booking data access
type BookingRepository struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

// NewBookingRepository creates a new booking repository
func NewBookingRepository(pool *pgxpool.Pool, queries *sqlc.Queries) *BookingRepository {
	return &BookingRepository{
		pool:    pool,
		queries: queries,
	}
}

// =============================================================================
// Main Booking Operations
// =============================================================================

// Create creates a new booking with slots and addons in a transaction
func (r *BookingRepository) Create(ctx context.Context, booking *models.Booking) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	// Create booking
	params := sqlc.CreateBookingParams{
		ID:            booking.ID,
		PackageID:     booking.PackageID,
		CustomerName:  booking.CustomerName,
		CustomerEmail: booking.CustomerEmail,
		CustomerPhone: booking.CustomerPhone,
		CustomerNotes: StringPtr(booking.CustomerNotes),
		Status:        string(booking.Status),
		TotalPrice:    DecimalToNumeric(booking.TotalAmount),
	}

	created, err := qtx.CreateBooking(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create booking: %w", err)
	}

	// Create slots
	for _, slot := range booking.Slots {
		slotTime, err := StringToTime(slot.Time)
		if err != nil {
			return fmt.Errorf("failed to parse slot time: %w", err)
		}

		slotParams := sqlc.CreateBookingSlotParams{
			BookingID: created.ID,
			Date:      TimeToDate(slot.Date),
			Time:      slotTime,
			ThemeID:   slot.ThemeID,
		}

		_, err = qtx.CreateBookingSlot(ctx, slotParams)
		if err != nil {
			return fmt.Errorf("failed to create booking slot: %w", err)
		}
	}

	// Create addons
	for _, addon := range booking.Addons {
		addonParams := sqlc.CreateBookingAddonParams{
			BookingID: created.ID,
			AddonID:   addon.AddonID,
			Quantity:  int32(addon.Quantity),
		}

		_, err = qtx.CreateBookingAddon(ctx, addonParams)
		if err != nil {
			return fmt.Errorf("failed to create booking addon: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Update booking with created values
	updatedBooking := r.toBookingModel(created)

	// Fetch the slots and addons that were just created with their generated IDs
	slotRows, err := r.queries.GetBookingSlots(ctx, created.ID)
	if err != nil {
		return fmt.Errorf("failed to get created booking slots: %w", err)
	}
	updatedBooking.Slots = r.toSlotModels(slotRows)

	addonRows, err := r.queries.GetBookingAddons(ctx, created.ID)
	if err != nil {
		return fmt.Errorf("failed to get created booking addons: %w", err)
	}
	updatedBooking.Addons = r.toAddonModels(addonRows)

	*booking = *updatedBooking
	return nil
}

// GetByID retrieves a booking with its slots and addons
func (r *BookingRepository) GetByID(ctx context.Context, id string) (*models.Booking, error) {
	// Get booking
	bookingRow, err := r.queries.GetBookingByID(ctx, id)
	if err != nil {
		return nil, err
	}

	booking := r.toBookingModel(bookingRow)

	// Get slots
	slotRows, err := r.queries.GetBookingSlots(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get booking slots: %w", err)
	}
	booking.Slots = r.toSlotModels(slotRows)

	// Get addons
	addonRows, err := r.queries.GetBookingAddons(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get booking addons: %w", err)
	}
	booking.Addons = r.toAddonModels(addonRows)

	return booking, nil
}

// ListAll retrieves all bookings with pagination
func (r *BookingRepository) ListAll(ctx context.Context, limit, offset int32) ([]*models.Booking, error) {
	params := sqlc.ListBookingsParams{
		Limit:  limit,
		Offset: offset,
	}

	rows, err := r.queries.ListBookings(ctx, params)
	if err != nil {
		return nil, err
	}

	bookings := make([]*models.Booking, len(rows))
	for i, row := range rows {
		bookings[i] = r.toBookingModel(row)
	}

	return bookings, nil
}

// GetByStatus retrieves bookings by status
func (r *BookingRepository) GetByStatus(ctx context.Context, status models.BookingStatus) ([]*models.Booking, error) {
	rows, err := r.queries.GetBookingsByStatus(ctx, string(status))
	if err != nil {
		return nil, err
	}

	bookings := make([]*models.Booking, len(rows))
	for i, row := range rows {
		bookings[i] = r.toBookingModel(row)
	}

	return bookings, nil
}

// GetByCustomerEmail retrieves bookings by customer email
func (r *BookingRepository) GetByCustomerEmail(ctx context.Context, email string) ([]*models.Booking, error) {
	rows, err := r.queries.GetBookingsByCustomerEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	bookings := make([]*models.Booking, len(rows))
	for i, row := range rows {
		bookings[i] = r.toBookingModel(row)
	}

	return bookings, nil
}

// UpdateStatus updates the status of a booking
func (r *BookingRepository) UpdateStatus(ctx context.Context, id string, status models.BookingStatus, adminNotes string) error {
	params := sqlc.UpdateBookingStatusParams{
		Column1:    id,
		Status:     string(status),
		AdminNotes: StringPtr(adminNotes),
	}

	_, err := r.queries.UpdateBookingStatus(ctx, params)
	return err
}

// UpdatePaymentScreenshot updates the payment screenshot URL
func (r *BookingRepository) UpdatePaymentScreenshot(ctx context.Context, id string, url string) error {
	params := sqlc.UpdateBookingPaymentScreenshotParams{
		Column1:              id,
		PaymentScreenshotUrl: StringPtr(url),
	}

	_, err := r.queries.UpdateBookingPaymentScreenshot(ctx, params)
	return err
}

// Count returns the total number of bookings
func (r *BookingRepository) Count(ctx context.Context) (int64, error) {
	return r.queries.CountBookings(ctx)
}

// =============================================================================
// Slot Operations
// =============================================================================

// GetBookingSlots retrieves slots for a booking
func (r *BookingRepository) GetBookingSlots(ctx context.Context, bookingID string) ([]models.BookingSlot, error) {
	rows, err := r.queries.GetBookingSlots(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	return r.toSlotModels(rows), nil
}

// GetBookedSlotsForThemeAndDate retrieves booked slots for a theme on a specific date
func (r *BookingRepository) GetBookedSlotsForThemeAndDate(ctx context.Context, themeID string, date string) ([]models.BookingSlot, error) {
	// Parse date string to time.Time
	// Note: Assuming date is in YYYY-MM-DD format
	dateTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	params := sqlc.GetBookedSlotsForThemeAndDateParams{
		Column1: themeID,
		Date:    TimeToDate(dateTime),
	}

	rows, err := r.queries.GetBookedSlotsForThemeAndDate(ctx, params)
	if err != nil {
		return nil, err
	}

	return r.toSlotModels(rows), nil
}

// GetBookedSlotsForAllThemesAndDate retrieves booked slots for all themes on a specific date
func (r *BookingRepository) GetBookedSlotsForAllThemesAndDate(ctx context.Context, date string) ([]models.BookingSlot, error) {
	// Parse date string to time.Time
	// Note: Assuming date is in YYYY-MM-DD format
	dateTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	rows, err := r.queries.GetBookedSlotsForAllThemesAndDate(ctx, TimeToDate(dateTime))
	if err != nil {
		return nil, err
	}

	return r.toSlotModels(rows), nil
}

// =============================================================================
// Addon Operations
// =============================================================================

// GetBookingAddons retrieves addons for a booking
func (r *BookingRepository) GetBookingAddons(ctx context.Context, bookingID string) ([]models.BookingAddon, error) {
	rows, err := r.queries.GetBookingAddons(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	return r.toAddonModels(rows), nil
}

// =============================================================================
// Conversion Helpers
// =============================================================================

// toBookingModel converts sqlc.Booking to models.Booking
func (r *BookingRepository) toBookingModel(row sqlc.Booking) *models.Booking {
	return &models.Booking{
		ID:                   row.ID,
		PackageID:            row.PackageID,
		CustomerName:         row.CustomerName,
		CustomerEmail:        row.CustomerEmail,
		CustomerPhone:        row.CustomerPhone,
		CustomerNotes:        StringVal(row.CustomerNotes),
		PaymentScreenshotURL: row.PaymentScreenshotUrl,
		Status:               models.BookingStatus(row.Status),
		TotalAmount:          NumericToDecimal(row.TotalPrice),
		AdminNotes:           StringVal(row.AdminNotes),
		CreatedAt:            TimestamptzToTime(row.CreatedAt),
		UpdatedAt:            TimestamptzToTime(row.UpdatedAt),
	}
}

// toSlotModels converts []sqlc.BookingSlot to []models.BookingSlot
func (r *BookingRepository) toSlotModels(rows []sqlc.BookingSlot) []models.BookingSlot {
	slots := make([]models.BookingSlot, len(rows))
	for i, row := range rows {
		slots[i] = models.BookingSlot{
			ID:        row.ID,
			BookingID: row.BookingID,
			ThemeID:   row.ThemeID,
			Date:      DateToTime(row.Date),
			Time:      TimeToString(row.Time),
		}
	}
	return slots
}

// toAddonModels converts []sqlc.BookingAddon to []models.BookingAddon
func (r *BookingRepository) toAddonModels(rows []sqlc.BookingAddon) []models.BookingAddon {
	addons := make([]models.BookingAddon, len(rows))
	for i, row := range rows {
		addons[i] = models.BookingAddon{
			ID:        row.ID,
			BookingID: row.BookingID,
			AddonID:   row.AddonID,
			Quantity:  int(row.Quantity),
		}
	}
	return addons
}
