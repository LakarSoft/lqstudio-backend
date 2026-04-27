package dto

import (
	"testing"

	"lqstudio-backend/internal/models"

	"github.com/shopspring/decimal"
)

func TestCreatePackageRequestToPackageModelIncludesModule(t *testing.T) {
	req := &CreatePackageRequest{
		ID:              "pkg-convo",
		Module:          "convocation",
		Name:            "Convocation Package",
		Description:     "For convocation sessions",
		DurationMinutes: 20,
		Price:           150,
		Discount:        5,
		Offers:          []string{"Photo session"},
		ImageURL:        "/uploads/packages/convo.jpg",
	}

	pkg := req.ToPackageModel()

	if pkg.Module != "convocation" {
		t.Fatalf("expected module convocation, got %q", pkg.Module)
	}
}

func TestToPackageResponseIncludesModule(t *testing.T) {
	pkg := &models.Package{
		ID:              "pkg-raya",
		Module:          "raya",
		Name:            "Raya Package",
		Description:     "Seasonal package",
		DurationMinutes: 20,
		Price:           decimal.NewFromInt(100),
		Discount:        decimal.NewFromInt(10),
		Offers:          []string{"Soft copy"},
		ImageURL:        "/uploads/packages/raya.jpg",
		IsActive:        true,
	}

	resp := ToPackageResponse(pkg)

	if resp.Module != "raya" {
		t.Fatalf("expected module raya, got %q", resp.Module)
	}
}
