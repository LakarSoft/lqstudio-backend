package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	stderr "errors"
	"fmt"
	"lqstudio-backend/internal/dto"
	"lqstudio-backend/internal/models"
	"lqstudio-backend/internal/repositories"
	"lqstudio-backend/pkg/email"
	"lqstudio-backend/pkg/errors"
	"lqstudio-backend/pkg/metrics"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// BookingService handles all booking-related business logic
type BookingService struct {
	packageRepo *repositories.PackageRepository
	themeRepo   *repositories.ThemeRepository
	addonRepo   *repositories.AddonRepository
	bookingRepo *repositories.BookingRepository
	emailClient *email.Client
	logger      *zap.Logger
}

// NewBookingService creates a new booking service
func NewBookingService(
	packageRepo *repositories.PackageRepository,
	themeRepo *repositories.ThemeRepository,
	addonRepo *repositories.AddonRepository,
	bookingRepo *repositories.BookingRepository,
	emailClient *email.Client,
	logger *zap.Logger,
) *BookingService {
	return &BookingService{
		packageRepo: packageRepo,
		themeRepo:   themeRepo,
		addonRepo:   addonRepo,
		bookingRepo: bookingRepo,
		emailClient: emailClient,
		logger:      logger,
	}
}

// CreateBooking creates a new booking with availability checking
func (s *BookingService) CreateBooking(ctx context.Context, req *dto.BookingRequest) (*dto.BookingResponse, error) {
	// 1. Validate and get package
	pkg, err := s.packageRepo.GetByID(ctx, req.PackageID)
	if err != nil {
		if stderr.Is(err, pgx.ErrNoRows) {
			return nil, errors.NewPackageNotFoundError(req.PackageID)
		}
		return nil, errors.NewDatabaseError("get package", err)
	}

	// Validate package is active
	if !pkg.IsActive {
		return nil, errors.NewBadRequestError(fmt.Sprintf("Package '%s' is not active", req.PackageID))
	}

	// 2. Validate slot count matches package requirements
	requiredSlots := pkg.RequiredSlots()
	if len(req.Slots) != requiredSlots {
		return nil, errors.NewInvalidSlotCountError(requiredSlots, len(req.Slots))
	}

	// 3. Validate all themes exist and are active
	themeIDs := make([]string, len(req.Slots))
	for i, slot := range req.Slots {
		themeIDs[i] = slot.ThemeID
	}

	// Get unique theme IDs
	uniqueThemeIDs := uniqueStrings(themeIDs)
	for _, themeID := range uniqueThemeIDs {
		theme, err := s.themeRepo.GetByID(ctx, themeID)
		if err != nil {
			if stderr.Is(err, pgx.ErrNoRows) {
				return nil, errors.NewThemeNotFoundError(themeID)
			}
			return nil, errors.NewDatabaseError(fmt.Sprintf("get theme %s", themeID), err)
		}
		if !theme.IsActive {
			return nil, errors.NewBadRequestError(fmt.Sprintf("Theme '%s' is not active", themeID))
		}
	}

	// 4. Validate all addons exist and are active
	if len(req.Addons) > 0 {
		addonIDs := make([]string, len(req.Addons))
		for i, addon := range req.Addons {
			addonIDs[i] = addon.AddonID
		}

		addons, err := s.addonRepo.GetByIDs(ctx, addonIDs)
		if err != nil {
			return nil, errors.NewDatabaseError("get addons", err)
		}

		// Validate all addons found and are active
		addonMap := make(map[string]*models.AddOn)
		for _, addon := range addons {
			addonMap[addon.ID] = addon
		}

		for _, reqAddon := range req.Addons {
			addon, exists := addonMap[reqAddon.AddonID]
			if !exists {
				return nil, errors.NewAddonNotFoundError(reqAddon.AddonID)
			}
			if !addon.IsActive {
				return nil, errors.NewBadRequestError(fmt.Sprintf("Addon '%s' is not active", reqAddon.AddonID))
			}
		}
	}

	// 5. Check slot availability for each (theme_id, date, time) combination
	for _, slot := range req.Slots {
		bookedSlots, err := s.bookingRepo.GetBookedSlotsForThemeAndDate(ctx, slot.ThemeID, slot.Date)
		if err != nil {
			return nil, errors.NewDatabaseError(fmt.Sprintf("check availability for theme %s on %s", slot.ThemeID, slot.Date), err)
		}

		// Check if this specific time is already booked
		for _, booked := range bookedSlots {
			if booked.Time == slot.Time {
				return nil, errors.NewSlotUnavailableError(slot.Date, slot.Time, slot.ThemeID)
			}
		}
	}

	// 6. Calculate prices
	// Package amount = package.FinalPrice() (which already applies discount)
	packageAmount := pkg.FinalPrice()

	// Addons amount = sum of (addon.Price * quantity) for each selected addon
	addonsAmount := decimal.Zero
	if len(req.Addons) > 0 {
		addonIDs := make([]string, len(req.Addons))
		for i, addon := range req.Addons {
			addonIDs[i] = addon.AddonID
		}

		addons, err := s.addonRepo.GetByIDs(ctx, addonIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to get addons for price calculation: %w", err)
		}

		// Build addon price map
		addonPriceMap := make(map[string]decimal.Decimal)
		for _, addon := range addons {
			addonPriceMap[addon.ID] = addon.Price
		}

		// Calculate addons total
		for _, reqAddon := range req.Addons {
			price, exists := addonPriceMap[reqAddon.AddonID]
			if !exists {
				continue // Already validated above
			}
			itemTotal := price.Mul(decimal.NewFromInt(int64(reqAddon.Quantity)))
			addonsAmount = addonsAmount.Add(itemTotal)
		}
	}

	// 7. Generate unique booking ID
	bookingID := generateBookingID()

	// 8. Create booking using DTO helper
	booking := req.ToBookingModel(packageAmount, addonsAmount)
	booking.ID = bookingID

	// 9. Create booking in repository
	if err := s.bookingRepo.Create(ctx, booking); err != nil {
		return nil, errors.NewDatabaseError("create booking", err)
	}

	// Increment booking counter metric
	metrics.BookingsTotal.Inc()

	// 10. Send email notifications (non-blocking - don't fail booking if emails fail)
	go s.sendBookingEmails(context.Background(), booking, pkg)

	// 11. Return booking response
	return dto.ToBookingResponse(booking), nil
}

// GetBookingByID retrieves a booking by ID
func (s *BookingService) GetBookingByID(ctx context.Context, bookingID string) (*dto.BookingResponse, error) {
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		if stderr.Is(err, pgx.ErrNoRows) {
			return nil, errors.NewBookingNotFoundError(bookingID)
		}
		return nil, errors.NewDatabaseError("get booking", err)
	}

	return dto.ToBookingResponse(booking), nil
}

// ListBookings lists bookings with optional filters and pagination
func (s *BookingService) ListBookings(ctx context.Context, limit, offset int32, statusFilter *models.BookingStatus, emailFilter *string) ([]*dto.BookingResponse, error) {
	var bookings []*models.Booking
	var err error

	// Apply filters
	if statusFilter != nil {
		bookings, err = s.bookingRepo.GetByStatus(ctx, *statusFilter)
	} else if emailFilter != nil {
		bookings, err = s.bookingRepo.GetByCustomerEmail(ctx, *emailFilter)
	} else {
		bookings, err = s.bookingRepo.ListAll(ctx, limit, offset)
	}

	if err != nil {
		return nil, errors.NewDatabaseError("list bookings", err)
	}

	// Apply manual pagination if using filtered results
	if statusFilter != nil || emailFilter != nil {
		start := int(offset)
		end := start + int(limit)
		if start > len(bookings) {
			start = len(bookings)
		}
		if end > len(bookings) {
			end = len(bookings)
		}
		bookings = bookings[start:end]
	}

	return dto.ToBookingsResponse(bookings), nil
}

// UpdateBookingStatus updates booking status
func (s *BookingService) UpdateBookingStatus(ctx context.Context, bookingID string, req *dto.UpdateBookingStatusRequest) (*dto.BookingResponse, error) {
	// Validate status is one of: PENDING, APPROVED, REJECTED
	status := models.BookingStatus(req.Status)
	if status != models.BookingStatusPending &&
		status != models.BookingStatusApproved &&
		status != models.BookingStatusRejected {
		return nil, errors.NewValidationError("Invalid booking status. Must be PENDING, APPROVED, or REJECTED")
	}

	// Check booking exists
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		if stderr.Is(err, pgx.ErrNoRows) {
			return nil, errors.NewBookingNotFoundError(bookingID)
		}
		return nil, errors.NewDatabaseError("get booking", err)
	}

	// Update status
	if err := s.bookingRepo.UpdateStatus(ctx, bookingID, status, req.AdminNotes); err != nil {
		return nil, errors.NewDatabaseError("update booking status", err)
	}

	// Get updated booking
	booking, err = s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return nil, errors.NewDatabaseError("get updated booking", err)
	}

	return dto.ToBookingResponse(booking), nil
}

// UpdatePaymentScreenshot updates the payment screenshot URL
func (s *BookingService) UpdatePaymentScreenshot(ctx context.Context, bookingID string, url string) (*dto.BookingResponse, error) {
	// Check booking exists
	_, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		if stderr.Is(err, pgx.ErrNoRows) {
			return nil, errors.NewBookingNotFoundError(bookingID)
		}
		return nil, errors.NewDatabaseError("get booking", err)
	}

	// Update payment screenshot
	if err := s.bookingRepo.UpdatePaymentScreenshot(ctx, bookingID, url); err != nil {
		return nil, errors.NewDatabaseError("update payment screenshot", err)
	}

	// Get updated booking
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return nil, errors.NewDatabaseError("get updated booking", err)
	}

	return dto.ToBookingResponse(booking), nil
}

// GetAvailability checks available time slots for a theme on a date
// If themeID is "all", it returns aggregated availability across all active themes
func (s *BookingService) GetAvailability(ctx context.Context, req *dto.AvailabilityRequest) (*dto.AvailabilityResponse, error) {
	// Optional: Validate package exists if packageID is provided
	if req.PackageID != "" {
		_, err := s.packageRepo.GetByID(ctx, req.PackageID)
		if err != nil {
			if stderr.Is(err, pgx.ErrNoRows) {
				return nil, errors.NewPackageNotFoundError(req.PackageID)
			}
			return nil, errors.NewDatabaseError("get package", err)
		}
	}

	// Check if requesting availability for all themes
	if req.ThemeID == "all" {
		return s.getAvailabilityForAllThemes(ctx, req)
	}

	// Single theme availability (existing logic)
	// Validate theme exists
	_, err := s.themeRepo.GetByID(ctx, req.ThemeID)
	if err != nil {
		if stderr.Is(err, pgx.ErrNoRows) {
			return nil, errors.NewThemeNotFoundError(req.ThemeID)
		}
		return nil, errors.NewDatabaseError("get theme", err)
	}

	// Generate time slots from 10:00 to 18:00 in 20-minute intervals
	allSlots := generateTimeSlots()

	// Get booked slots for this theme on this date
	bookedSlots, err := s.bookingRepo.GetBookedSlotsForThemeAndDate(ctx, req.ThemeID, req.Date)
	if err != nil {
		return nil, errors.NewDatabaseError("get booked slots", err)
	}

	// Build a map of booked times for quick lookup
	bookedTimes := make(map[string]bool)
	for _, slot := range bookedSlots {
		bookedTimes[slot.Time] = true
	}

	// Build availability response
	availabilitySlots := make([]dto.AvailableSlotInfo, len(allSlots))
	for i, slotTime := range allSlots {
		availabilitySlots[i] = dto.AvailableSlotInfo{
			Time:      slotTime,
			Available: !bookedTimes[slotTime],
			ThemeID:   req.ThemeID,
		}
	}

	return &dto.AvailabilityResponse{
		Date:  req.Date,
		Slots: availabilitySlots,
	}, nil
}

// getAvailabilityForAllThemes checks availability across all active themes
// A time slot is available ONLY if ALL themes are available (no themes booked)
// If even one theme is booked at a time slot, it shows as unavailable
func (s *BookingService) getAvailabilityForAllThemes(ctx context.Context, req *dto.AvailabilityRequest) (*dto.AvailabilityResponse, error) {
	// Get all active themes
	activeThemes, err := s.themeRepo.GetActive(ctx)
	if err != nil {
		return nil, errors.NewDatabaseError("get active themes", err)
	}

	// If no active themes, all slots are unavailable
	if len(activeThemes) == 0 {
		allSlots := generateTimeSlots()
		availabilitySlots := make([]dto.AvailableSlotInfo, len(allSlots))
		for i, slotTime := range allSlots {
			availabilitySlots[i] = dto.AvailableSlotInfo{
				Time:      slotTime,
				Available: false,
			}
		}
		return &dto.AvailabilityResponse{
			Date:  req.Date,
			Slots: availabilitySlots,
		}, nil
	}

	// Generate time slots from 10:00 to 18:00 in 20-minute intervals
	allSlots := generateTimeSlots()

	// Get all booked slots for all themes on this date
	bookedSlots, err := s.bookingRepo.GetBookedSlotsForAllThemesAndDate(ctx, req.Date)
	if err != nil {
		return nil, errors.NewDatabaseError("get booked slots for all themes", err)
	}

	// Build a map: time -> set of booked theme IDs
	bookedThemesByTime := make(map[string]map[string]bool)
	for _, slot := range bookedSlots {
		if bookedThemesByTime[slot.Time] == nil {
			bookedThemesByTime[slot.Time] = make(map[string]bool)
		}
		bookedThemesByTime[slot.Time][slot.ThemeID] = true
	}

	// Calculate availability: available ONLY if NO themes are booked at that time
	availabilitySlots := make([]dto.AvailableSlotInfo, len(allSlots))
	for i, slotTime := range allSlots {
		bookedCount := len(bookedThemesByTime[slotTime])
		// Available only if no themes are booked (i.e., all themes are free)
		availabilitySlots[i] = dto.AvailableSlotInfo{
			Time:      slotTime,
			Available: bookedCount == 0,
			// ThemeID is omitted (will be empty due to omitempty tag)
		}
	}

	return &dto.AvailabilityResponse{
		Date:  req.Date,
		Slots: availabilitySlots,
	}, nil
}

// =============================================================================
// Helper Functions
// =============================================================================

// sendBookingEmails sends confirmation and notification emails for a booking
func (s *BookingService) sendBookingEmails(ctx context.Context, booking *models.Booking, pkg *models.Package) {
	// Skip if email client is not configured
	if s.emailClient == nil {
		s.logger.Info("Email client not configured, skipping email notifications",
			zap.String("booking_id", booking.ID),
		)
		return
	}

	// Prepare slot information for email
	slots := make([]email.SlotInfo, len(booking.Slots))
	for i, slot := range booking.Slots {
		// Get theme name
		theme, err := s.themeRepo.GetByID(ctx, slot.ThemeID)
		themeName := slot.ThemeID // fallback to ID if theme not found
		if err == nil {
			themeName = theme.Name
		}

		slots[i] = email.SlotInfo{
			ThemeName: themeName,
			Date:      slot.Date.Format("Monday, 02 Jan 2006"),
			Time:      slot.Time,
		}
	}

	// Prepare addon information for email
	addons := make([]email.AddonInfo, len(booking.Addons))
	for i, bookingAddon := range booking.Addons {
		// Get addon details
		addon, err := s.addonRepo.GetByID(ctx, bookingAddon.AddonID)
		if err != nil {
			s.logger.Warn("Failed to get addon details for email",
				zap.String("booking_id", booking.ID),
				zap.String("addon_id", bookingAddon.AddonID),
				zap.Error(err),
			)
			continue
		}

		itemTotal := addon.Price.Mul(decimal.NewFromInt(int64(bookingAddon.Quantity)))
		addons[i] = email.AddonInfo{
			Name:     addon.Name,
			Quantity: bookingAddon.Quantity,
			Price:    itemTotal.StringFixed(2),
		}
	}

	// Send customer confirmation email
	if err := s.emailClient.SendBookingConfirmation(
		booking.CustomerEmail,
		booking,
		pkg.Name,
		slots,
		addons,
	); err != nil {
		s.logger.Error("Failed to send booking confirmation email",
			zap.String("booking_id", booking.ID),
			zap.Error(err),
		)
		metrics.EmailFailuresTotal.WithLabelValues("customer_confirmation").Inc()
		// Don't return - still try to send admin notification
	} else {
		metrics.EmailNotificationsTotal.WithLabelValues("customer_confirmation").Inc()
	}

	// Send admin notification email
	if err := s.emailClient.SendAdminNotification(
		booking,
		pkg.Name,
		slots,
		addons,
	); err != nil {
		s.logger.Error("Failed to send admin notification email",
			zap.String("booking_id", booking.ID),
			zap.Error(err),
		)
		metrics.EmailFailuresTotal.WithLabelValues("admin_notification").Inc()
	} else {
		metrics.EmailNotificationsTotal.WithLabelValues("admin_notification").Inc()
	}
}

// generateBookingID generates a unique booking ID in format: bkg-{timestamp}-{random}
func generateBookingID() string {
	timestamp := time.Now().Unix()
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	randomHex := hex.EncodeToString(randomBytes)
	return fmt.Sprintf("bkg-%d-%s", timestamp, randomHex)
}

// generateTimeSlots generates time slots from 10:00 to 18:00 in 20-minute intervals
func generateTimeSlots() []string {
	slots := []string{}
	hour := 10
	minute := 0

	for {
		// Stop at 18:00
		if hour >= 18 {
			break
		}

		// Format time as HH:MM
		timeStr := fmt.Sprintf("%02d:%02d", hour, minute)
		slots = append(slots, timeStr)

		// Add 20 minutes
		minute += 20
		if minute >= 60 {
			minute = 0
			hour++
		}
	}

	return slots
}

// uniqueStrings returns unique strings from a slice
func uniqueStrings(strs []string) []string {
	seen := make(map[string]bool)
	unique := []string{}

	for _, str := range strs {
		if !seen[str] {
			seen[str] = true
			unique = append(unique, str)
		}
	}

	return unique
}
