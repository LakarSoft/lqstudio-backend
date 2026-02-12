package dto

import (
	"lqstudio-backend/internal/models"

	"github.com/shopspring/decimal"
)

// AddonResponse represents add-on information for API responses
// Matches the frontend Addon entity structure with camelCase JSON tags
type AddonResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price"`          // Price per unit in MYR
	Unit        string  `json:"unit,omitempty"` // Unit of measurement (e.g., "pax", "pc", "person")
	IsActive    bool    `json:"isActive"`
}

// CreateAddonRequest for admin add-on creation
// Matches the frontend expectations with camelCase JSON tags
type CreateAddonRequest struct {
	ID          string  `json:"id" validate:"required,max=50"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,min=0"`
	Unit        string  `json:"unit"`
}

// UpdateAddonRequest for admin add-on updates
// Matches the frontend expectations with camelCase JSON tags
type UpdateAddonRequest struct {
	ID          string  `json:"id" validate:"required,max=50"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,min=0"`
	Unit        string  `json:"unit"`
	IsActive    bool    `json:"isActive"`
}

// SelectedAddon represents an add-on selection in a booking request
// Matches the frontend SelectedAddon structure with camelCase JSON tags
type SelectedAddon struct {
	AddonID  string `json:"addonId" validate:"required"`
	Quantity int    `json:"quantity" validate:"required,min=1"`
}

// ToAddonResponse converts a domain model AddOn to an AddonResponse DTO
// Handles decimal.Decimal to float64 conversion
func ToAddonResponse(addon *models.AddOn) *AddonResponse {
	if addon == nil {
		return nil
	}

	// Convert decimal.Decimal to float64
	price, _ := addon.Price.Float64()

	return &AddonResponse{
		ID:          addon.ID,
		Name:        addon.Name,
		Description: addon.Description,
		Price:       price,
		Unit:        addon.Unit,
		IsActive:    addon.IsActive,
	}
}

// ToAddonsResponse converts a slice of domain model AddOns to a slice of AddonResponse DTOs
func ToAddonsResponse(addons []*models.AddOn) []*AddonResponse {
	if addons == nil {
		return []*AddonResponse{}
	}

	responses := make([]*AddonResponse, len(addons))
	for i, addon := range addons {
		responses[i] = ToAddonResponse(addon)
	}
	return responses
}

// ToAddonModel converts CreateAddonRequest to domain model AddOn
// Handles float64 to decimal.Decimal conversion
func (r *CreateAddonRequest) ToAddonModel() *models.AddOn {
	return &models.AddOn{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Price:       decimal.NewFromFloat(r.Price),
		Unit:        r.Unit,
		IsActive:    true, // New addons are active by default
	}
}

// ToAddonModel converts UpdateAddonRequest to domain model AddOn
// Handles float64 to decimal.Decimal conversion
func (r *UpdateAddonRequest) ToAddonModel() *models.AddOn {
	return &models.AddOn{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Price:       decimal.NewFromFloat(r.Price),
		Unit:        r.Unit,
		IsActive:    r.IsActive,
	}
}
