package dto

import (
	"testing"

	"lqstudio-backend/internal/models"

	"github.com/shopspring/decimal"
)

func TestCreateThemeRequestToThemeModelIncludesModule(t *testing.T) {
	req := &CreateThemeRequest{
		ID:          "theme-convo",
		Module:      "convocation",
		Name:        "Convocation Theme",
		Description: "Graduation theme",
		ImageURL:    "/uploads/themes/convo.jpg",
		Price:       50,
	}

	theme := req.ToThemeModel()

	if theme.Module != "convocation" {
		t.Fatalf("expected module convocation, got %q", theme.Module)
	}
}

func TestToThemeResponseIncludesModule(t *testing.T) {
	theme := &models.Theme{
		ID:          "theme-raya",
		Module:      "raya",
		Name:        "Raya Theme",
		Description: "Seasonal theme",
		ImageURL:    "/uploads/themes/raya.jpg",
		Price:       decimal.NewFromInt(25),
		IsActive:    true,
	}

	resp := ToThemeResponse(theme)

	if resp.Module != "raya" {
		t.Fatalf("expected module raya, got %q", resp.Module)
	}
}
