package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// Theme represents a photography theme
type Theme struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	ImageURL    string          `json:"imageUrl"`
	Price       decimal.Decimal `json:"price"`
	IsActive    bool            `json:"isActive"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}
