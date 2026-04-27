package dto

import (
	"testing"
	"time"

	"lqstudio-backend/internal/models"

	"github.com/shopspring/decimal"
)

func TestToBookingDetailsResponseIncludesNestedDetails(t *testing.T) {
	createdAt := time.Date(2026, 4, 25, 2, 30, 0, 0, time.UTC)
	updatedAt := createdAt.Add(5 * time.Minute)

	booking := &models.Booking{
		ID:            "booking_123",
		PackageID:     "pkg-raya-single",
		Status:        models.BookingStatusPending,
		CustomerName:  "Anas",
		CustomerEmail: "anas@example.com",
		CustomerPhone: "0123456789",
		CustomerNotes: "Bring props",
		TotalAmount:   decimal.NewFromInt(180),
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
		Slots: []models.BookingSlot{
			{
				ThemeID: "theme-a",
				Date:    time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
				Time:    "14:20",
			},
		},
		Addons: []models.BookingAddon{
			{
				AddonID:  "addon-1",
				Quantity: 2,
			},
		},
	}

	pkg := &models.Package{
		ID:              "pkg-raya-single",
		Module:          models.ModuleRaya,
		Name:            "Raya Single",
		Description:     "Single slot package",
		DurationMinutes: 20,
		Price:           decimal.NewFromInt(200),
		Discount:        decimal.NewFromInt(10),
		IsActive:        true,
	}

	theme := &models.Theme{
		ID:          "theme-a",
		Module:      models.ModuleRaya,
		Name:        "Theme A",
		Description: "Main setup",
		ImageURL:    "/images/theme-a.jpg",
		Price:       decimal.Zero,
		IsActive:    true,
	}

	addon := &models.AddOn{
		ID:          "addon-1",
		Module:      models.ModuleRaya,
		Name:        "Extra Edited Photos",
		Description: "More edited outputs",
		Price:       decimal.NewFromInt(15),
		Unit:        "pc",
		IsActive:    true,
	}

	resp := ToBookingDetailsResponse(
		booking,
		pkg,
		map[string]*models.Theme{"theme-a": theme},
		map[string]*models.AddOn{"addon-1": addon},
	)

	if resp == nil {
		t.Fatal("expected response, got nil")
	}
	if resp.Package == nil || resp.Package.ID != pkg.ID {
		t.Fatalf("expected package details for %q, got %+v", pkg.ID, resp.Package)
	}
	if len(resp.Slots) != 1 {
		t.Fatalf("expected 1 slot, got %d", len(resp.Slots))
	}
	if resp.Slots[0].Theme == nil || resp.Slots[0].Theme.ID != theme.ID {
		t.Fatalf("expected nested theme details for %q, got %+v", theme.ID, resp.Slots[0].Theme)
	}
	if len(resp.Addons) != 1 {
		t.Fatalf("expected 1 addon, got %d", len(resp.Addons))
	}
	if resp.Addons[0].Addon == nil || resp.Addons[0].Addon.ID != addon.ID {
		t.Fatalf("expected nested addon details for %q, got %+v", addon.ID, resp.Addons[0].Addon)
	}
	if resp.TotalPrice != 180 {
		t.Fatalf("expected total price 180, got %v", resp.TotalPrice)
	}
	if resp.CreatedAt != createdAt.Format(time.RFC3339) {
		t.Fatalf("expected createdAt %q, got %q", createdAt.Format(time.RFC3339), resp.CreatedAt)
	}
}
