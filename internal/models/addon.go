package models

import (
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// AddOn represents an optional add-on service
type AddOn struct {
	ID          string          `json:"id"`
	Module      string          `json:"module"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price"`
	Unit        string          `json:"unit"` // e.g., "person", "set", "session"
	IsActive    bool            `json:"isActive"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

// ValidateAddon checks if addon data is valid.
func (a *AddOn) ValidateAddon() error {
	a.Module = strings.TrimSpace(strings.ToLower(a.Module))
	if a.Module != "raya" && a.Module != "convocation" {
		return ErrInvalidAddonModule
	}
	return nil
}

var ErrInvalidAddonModule = NewValidationError("module must be either raya or convocation")
