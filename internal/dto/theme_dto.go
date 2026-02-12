package dto

import (
	"lqstudio-backend/internal/models"

	"github.com/shopspring/decimal"
)

// ThemeResponse represents theme information for API responses
// Matches the frontend Theme entity structure with camelCase JSON tags
type ThemeResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageURL    string  `json:"imageUrl"`
	Price       float64 `json:"price"` // Additional cost in MYR (currently 0 for all themes)
	IsActive    bool    `json:"isActive"`
}

// CreateThemeRequest for admin theme creation
// Matches the frontend expectations with camelCase JSON tags
type CreateThemeRequest struct {
	ID          string  `json:"id" validate:"required,max=50"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	ImageURL    string  `json:"imageUrl" validate:"required"`
	Price       float64 `json:"price,omitempty"` // Optional - defaults to 0 if not provided
}

// UpdateThemeRequest for admin theme updates
// Matches the frontend expectations with camelCase JSON tags
type UpdateThemeRequest struct {
	ID          string  `json:"id" validate:"required,max=50"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	ImageURL    string  `json:"imageUrl" validate:"required"`
	Price       float64 `json:"price,omitempty"` // Optional - defaults to 0 if not provided
	IsActive    bool    `json:"isActive"`
}

// ToThemeResponse converts a domain model Theme to a ThemeResponse DTO
// Handles decimal.Decimal to float64 conversion
func ToThemeResponse(theme *models.Theme) *ThemeResponse {
	if theme == nil {
		return nil
	}

	// Convert decimal.Decimal to float64
	price, _ := theme.Price.Float64()

	return &ThemeResponse{
		ID:          theme.ID,
		Name:        theme.Name,
		Description: theme.Description,
		ImageURL:    theme.ImageURL,
		Price:       price,
		IsActive:    theme.IsActive,
	}
}

// ToThemesResponse converts a slice of domain model Themes to a slice of ThemeResponse DTOs
func ToThemesResponse(themes []*models.Theme) []*ThemeResponse {
	if themes == nil {
		return []*ThemeResponse{}
	}

	responses := make([]*ThemeResponse, len(themes))
	for i, theme := range themes {
		responses[i] = ToThemeResponse(theme)
	}
	return responses
}

// ToThemeModel converts CreateThemeRequest to domain model Theme
// Handles float64 to decimal.Decimal conversion
func (r *CreateThemeRequest) ToThemeModel() *models.Theme {
	return &models.Theme{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		ImageURL:    r.ImageURL,
		Price:       decimal.NewFromFloat(r.Price),
		IsActive:    true, // New themes are active by default
	}
}

// ToThemeModel converts UpdateThemeRequest to domain model Theme
// Handles float64 to decimal.Decimal conversion
func (r *UpdateThemeRequest) ToThemeModel() *models.Theme {
	return &models.Theme{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		ImageURL:    r.ImageURL,
		Price:       decimal.NewFromFloat(r.Price),
		IsActive:    r.IsActive,
	}
}
