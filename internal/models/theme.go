package models

import (
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// Theme represents a photography theme
type Theme struct {
	ID          string          `json:"id"`
	Module      string          `json:"module"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	ImageURL    string          `json:"imageUrl"`
	Price       decimal.Decimal `json:"price"`
	IsActive    bool            `json:"isActive"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

// ValidateTheme checks if theme data is valid.
func (t *Theme) ValidateTheme() error {
	t.Module = strings.TrimSpace(strings.ToLower(t.Module))
	if t.Module != "raya" && t.Module != "convocation" {
		return ErrInvalidThemeModule
	}
	return nil
}

var ErrInvalidThemeModule = NewValidationError("module must be either raya or convocation")
