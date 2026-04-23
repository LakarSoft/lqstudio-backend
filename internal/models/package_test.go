package models

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestPackageValidatePackageRequiresModule(t *testing.T) {
	pkg := &Package{
		ID:              "pkg-1",
		Name:            "Test Package",
		DurationMinutes: 20,
		Price:           decimal.NewFromInt(100),
		Discount:        decimal.Zero,
	}

	err := pkg.ValidatePackage()
	if err != ErrInvalidModule {
		t.Fatalf("expected ErrInvalidModule, got %v", err)
	}
}

func TestPackageValidatePackageNormalizesModule(t *testing.T) {
	pkg := &Package{
		ID:              "pkg-1",
		Module:          "  Convocation  ",
		Name:            "Test Package",
		DurationMinutes: 30,
		Price:           decimal.NewFromInt(100),
		Discount:        decimal.Zero,
	}

	if err := pkg.ValidatePackage(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if pkg.Module != "convocation" {
		t.Fatalf("expected normalized module convocation, got %q", pkg.Module)
	}
}

func TestPackageRequiredSlotsUsesRayaInterval(t *testing.T) {
	pkg := &Package{
		Module:          ModuleRaya,
		DurationMinutes: 40,
	}

	if got := pkg.RequiredSlots(); got != 2 {
		t.Fatalf("expected 2 required slots, got %d", got)
	}
}

func TestPackageRequiredSlotsUsesConvocationInterval(t *testing.T) {
	pkg := &Package{
		Module:          ModuleConvocation,
		DurationMinutes: 30,
	}

	if got := pkg.RequiredSlots(); got != 1 {
		t.Fatalf("expected 1 required slot, got %d", got)
	}
}

func TestPackageValidatePackageRejectsInvalidConvocationDuration(t *testing.T) {
	pkg := &Package{
		ID:              "pkg-1",
		Module:          ModuleConvocation,
		Name:            "Test Package",
		DurationMinutes: 20,
		Price:           decimal.NewFromInt(100),
		Discount:        decimal.Zero,
	}

	if err := pkg.ValidatePackage(); err == nil {
		t.Fatal("expected validation error for non-30-minute convocation duration")
	}
}
