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
		DurationMinutes: 20,
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
