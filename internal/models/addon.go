package models

import (
	"time"

	"github.com/shopspring/decimal"
)

// AddOn represents an optional add-on service
type AddOn struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price"`
	Unit        string          `json:"unit"` // e.g., "person", "set", "session"
	IsActive    bool            `json:"isActive"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}
