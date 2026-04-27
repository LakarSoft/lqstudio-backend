package dto

import (
	"testing"

	"lqstudio-backend/internal/models"

	"github.com/shopspring/decimal"
)

func TestCreateAddonRequestToAddonModelIncludesModule(t *testing.T) {
	req := &CreateAddonRequest{
		ID:          "addon-convo",
		Module:      "convocation",
		Name:        "Frame",
		Description: "Convocation frame",
		Price:       99,
		Unit:        "pc",
	}

	addon := req.ToAddonModel()

	if addon.Module != "convocation" {
		t.Fatalf("expected module convocation, got %q", addon.Module)
	}
}

func TestToAddonResponseIncludesModule(t *testing.T) {
	addon := &models.AddOn{
		ID:          "addon-raya",
		Module:      "raya",
		Name:        "Album",
		Description: "Seasonal album",
		Price:       decimal.NewFromInt(100),
		Unit:        "pc",
		IsActive:    true,
	}

	resp := ToAddonResponse(addon)

	if resp.Module != "raya" {
		t.Fatalf("expected module raya, got %q", resp.Module)
	}
}
