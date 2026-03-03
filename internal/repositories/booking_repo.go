package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"lqstudio-backend/internal/database/sqlc"
	"lqstudio-backend/internal/dto"
	"lqstudio-backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

// UpdateAdminNotes updates only the admin_notes field of a booking
func (r *BookingRepository) UpdateAdminNotes(ctx context.Context, id string, adminNotes string) error {
	params := sqlc.UpdateBookingAdminNotesParams{
		Column1:    id,
		AdminNotes: StringPtr(adminNotes),
	}

	_, err := r.queries.UpdateBookingAdminNotes(ctx, params)
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

// UpdateBooking replaces slots and addons for an existing booking in a single transaction.
// It also updates customer info and recalculates total_price.
func (r *BookingRepository) UpdateBooking(ctx context.Context, bookingID string, booking *models.Booking) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	// Delete existing slots and addons
	if err := qtx.DeleteBookingSlots(ctx, bookingID); err != nil {
		return fmt.Errorf("failed to delete booking slots: %w", err)
	}
	if err := qtx.DeleteBookingAddons(ctx, bookingID); err != nil {
		return fmt.Errorf("failed to delete booking addons: %w", err)
	}

	// Insert new slots
	for _, slot := range booking.Slots {
		slotTime, err := StringToTime(slot.Time)
		if err != nil {
			return fmt.Errorf("failed to parse slot time: %w", err)
		}
		_, err = qtx.CreateBookingSlot(ctx, sqlc.CreateBookingSlotParams{
			BookingID: bookingID,
			Date:      TimeToDate(slot.Date),
			Time:      slotTime,
			ThemeID:   slot.ThemeID,
		})
		if err != nil {
			return fmt.Errorf("failed to create booking slot: %w", err)
		}
	}

	// Insert new addons
	for _, addon := range booking.Addons {
		_, err = qtx.CreateBookingAddon(ctx, sqlc.CreateBookingAddonParams{
			BookingID: bookingID,
			AddonID:   addon.AddonID,
			Quantity:  int32(addon.Quantity),
		})
		if err != nil {
			return fmt.Errorf("failed to create booking addon: %w", err)
		}
	}

	// Update booking core fields
	updated, err := qtx.UpdateBookingDetails(ctx, sqlc.UpdateBookingDetailsParams{
		Column1:       bookingID,
		CustomerName:  booking.CustomerName,
		CustomerEmail: booking.CustomerEmail,
		CustomerPhone: booking.CustomerPhone,
		CustomerNotes: StringPtr(booking.CustomerNotes),
		TotalPrice:    DecimalToNumeric(booking.TotalAmount),
	})
	if err != nil {
		return fmt.Errorf("failed to update booking details: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Reflect committed state back into the passed model
	updatedModel := r.toBookingModel(updated)

	slotRows, err := r.queries.GetBookingSlots(ctx, bookingID)
	if err != nil {
		return fmt.Errorf("failed to fetch updated slots: %w", err)
	}
	updatedModel.Slots = r.toSlotModels(slotRows)

	addonRows, err := r.queries.GetBookingAddons(ctx, bookingID)
	if err != nil {
		return fmt.Errorf("failed to fetch updated addons: %w", err)
	}
	updatedModel.Addons = r.toAddonModels(addonRows)

	*booking = *updatedModel
	return nil
}

// Count returns the total number of bookings
func (r *BookingRepository) Count(ctx context.Context) (int64, error) {
	return r.queries.CountBookings(ctx)
}

// ListWithFilters retrieves bookings with dynamic filtering, sorting, and pagination.
// Uses raw pgxpool query because sqlc does not support dynamic ORDER BY or complex conditional WHERE.
func (r *BookingRepository) ListWithFilters(ctx context.Context, filters *dto.BookingFilters) ([]*models.Booking, int, error) {
	var conditions []string
	var args []interface{}
	argIdx := 1

	if filters.Status != "" {
		conditions = append(conditions, fmt.Sprintf("b.status = $%d", argIdx))
		args = append(args, filters.Status)
		argIdx++
	}

	if filters.Email != "" {
		conditions = append(conditions, fmt.Sprintf("b.customer_email ILIKE $%d", argIdx))
		args = append(args, "%"+filters.Email+"%")
		argIdx++
	}

	if filters.PackageID != "" {
		conditions = append(conditions, fmt.Sprintf("b.package_id = $%d", argIdx))
		args = append(args, filters.PackageID)
		argIdx++
	}

	if filters.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(b.customer_name ILIKE $%d OR b.customer_email ILIKE $%d OR b.customer_phone ILIKE $%d)", argIdx, argIdx, argIdx))
		args = append(args, "%"+filters.Search+"%")
		argIdx++
	}

	needsSlotJoin := filters.ThemeID != "" || filters.SlotDate != "" || filters.DateFrom != "" || filters.DateTo != ""

	if filters.DateFrom != "" {
		conditions = append(conditions, fmt.Sprintf("bs.date >= $%d::date", argIdx))
		args = append(args, filters.DateFrom)
		argIdx++
	}

	if filters.DateTo != "" {
		conditions = append(conditions, fmt.Sprintf("bs.date <= $%d::date", argIdx))
		args = append(args, filters.DateTo)
		argIdx++
	}

	if filters.ThemeID != "" {
		conditions = append(conditions, fmt.Sprintf("bs.theme_id = $%d", argIdx))
		args = append(args, filters.ThemeID)
		argIdx++
	}

	if filters.SlotDate != "" {
		conditions = append(conditions, fmt.Sprintf("bs.date = $%d::date", argIdx))
		args = append(args, filters.SlotDate)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	joinClause := ""
	if needsSlotJoin {
		joinClause = "LEFT JOIN booking_slots bs ON b.id = bs.booking_id"
	}

	// Map sort field to SQL column name (whitelist to prevent injection)
	orderColumn := "b.created_at"
	subOrderColumn := "created_at"
	switch filters.SortBy {
	case "createdAt":
		orderColumn = "b.created_at"
		subOrderColumn = "created_at"
	case "updatedAt":
		orderColumn = "b.updated_at"
		subOrderColumn = "updated_at"
	case "totalPrice":
		orderColumn = "b.total_price"
		subOrderColumn = "total_price"
	case "status":
		orderColumn = "b.status"
		subOrderColumn = "status"
	}

	orderDir := "DESC"
	if strings.ToLower(filters.Order) == "asc" {
		orderDir = "ASC"
	}

	// Count distinct bookings (handles potential JOIN duplicates)
	countQuery := fmt.Sprintf(
		"SELECT COUNT(DISTINCT b.id) FROM bookings b %s %s",
		joinClause, whereClause,
	)

	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count bookings: %w", err)
	}

	// Data query: use DISTINCT ON when joining slots to avoid duplicate rows
	const selectFields = `b.id, b.package_id, b.customer_name, b.customer_email, b.customer_phone,
		b.customer_notes, b.payment_screenshot_url, b.status, b.total_price, b.admin_notes,
		b.created_at, b.updated_at`

	var dataQuery string
	if needsSlotJoin {
		// DISTINCT ON (b.id) deduplicate; wrap in subquery for correct outer ORDER BY
		innerQuery := fmt.Sprintf(
			"SELECT DISTINCT ON (b.id) %s FROM bookings b %s %s ORDER BY b.id, %s %s",
			selectFields, joinClause, whereClause, orderColumn, orderDir,
		)
		dataQuery = fmt.Sprintf(
			"SELECT * FROM (%s) sub ORDER BY %s %s LIMIT $%d OFFSET $%d",
			innerQuery, subOrderColumn, orderDir, argIdx, argIdx+1,
		)
	} else {
		dataQuery = fmt.Sprintf(
			"SELECT %s FROM bookings b %s ORDER BY %s %s LIMIT $%d OFFSET $%d",
			selectFields, whereClause, orderColumn, orderDir, argIdx, argIdx+1,
		)
	}
	args = append(args, filters.Limit, filters.Offset())

	rows, err := r.pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query bookings: %w", err)
	}
	defer rows.Close()

	var bookings []*models.Booking
	for rows.Next() {
		var b models.Booking
		var customerNotes, paymentURL, adminNotes *string
		var status string
		var totalPrice pgtype.Numeric
		var createdAt, updatedAt pgtype.Timestamptz

		if err := rows.Scan(
			&b.ID, &b.PackageID, &b.CustomerName, &b.CustomerEmail, &b.CustomerPhone,
			&customerNotes, &paymentURL, &status, &totalPrice, &adminNotes,
			&createdAt, &updatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan booking row: %w", err)
		}

		b.CustomerNotes = StringVal(customerNotes)
		b.PaymentScreenshotURL = paymentURL
		b.Status = models.BookingStatus(status)
		b.TotalAmount = NumericToDecimal(totalPrice)
		b.AdminNotes = StringVal(adminNotes)
		b.CreatedAt = TimestamptzToTime(createdAt)
		b.UpdatedAt = TimestamptzToTime(updatedAt)

		bookings = append(bookings, &b)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating booking rows: %w", err)
	}

	if bookings == nil {
		bookings = []*models.Booking{}
	}

	return bookings, total, nil
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
